const API_BASE = import.meta.env.VITE_API_URL || '';

export interface AuthData {
  accessToken: string;
  refreshToken: string;
  username: string;
  email: string;
  isNewUser: boolean;
}

export interface UserProfile {
  username: string;
  email: string;
  elo: number;
  xp: number;
  level: number;
  streak: number;
  rank: string;
  completed_ch: number;
  created_at: string;
}

export interface ChapterProgress {
  number: number;
  title: string;
  city: string;
  difficulty: string;
  status: 'locked' | 'active' | 'complete';
}

class ApiService {
  private accessToken: string | null = localStorage.getItem('access_token');
  private refreshToken: string | null = localStorage.getItem('refresh_token');
  private username: string | null = localStorage.getItem('username');

  getApiBase(): string {
    return API_BASE;
  }

  getUsername(): string | null {
    return this.username;
  }

  isLoggedIn(): boolean {
    return !!this.accessToken;
  }

  setTokens(data: AuthData) {
    this.accessToken = data.accessToken;
    this.refreshToken = data.refreshToken;
    this.username = data.username;
    localStorage.setItem('access_token', data.accessToken);
    localStorage.setItem('refresh_token', data.refreshToken);
    localStorage.setItem('username', data.username);
  }

  logout() {
    this.accessToken = null;
    this.refreshToken = null;
    this.username = null;
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    localStorage.removeItem('username');
  }

  // Request wrapper with auto token refresh
  async request<T>(path: string, options: RequestInit = {}): Promise<T> {
    const headers = new Headers(options.headers || {});
    if (this.accessToken) {
      headers.set('Authorization', `Bearer ${this.accessToken}`);
    }
    headers.set('Content-Type', 'application/json');

    const url = `${API_BASE}${path}`;
    const response = await fetch(url, { ...options, headers });

    if (response.status === 401 && this.refreshToken) {
      // Access token expired, attempt refresh
      try {
        const refreshed = await this.refreshTokens();
        if (refreshed) {
          // Retry request with new token
          const retryHeaders = new Headers(options.headers || {});
          retryHeaders.set('Authorization', `Bearer ${this.accessToken}`);
          retryHeaders.set('Content-Type', 'application/json');
          const retryResponse = await fetch(url, { ...options, headers: retryHeaders });
          if (!retryResponse.ok) {
            throw new Error(`Request failed after token refresh: ${retryResponse.statusText}`);
          }
          return await retryResponse.json() as T;
        }
      } catch (err) {
        this.logout();
        throw new Error('Session expired. Please log in again.');
      }
    }

    if (!response.ok) {
      const errBody = await response.json().catch(() => ({}));
      throw new Error(errBody.error || response.statusText);
    }

    return await response.json() as T;
  }

  private async refreshTokens(): Promise<boolean> {
    if (!this.refreshToken) return false;

    const response = await fetch(`${API_BASE}/api/auth/refresh`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: this.refreshToken }),
    });

    if (!response.ok) {
      this.logout();
      return false;
    }

    const data = await response.json() as { access_token: string; refresh_token: string };
    this.accessToken = data.access_token;
    this.refreshToken = data.refresh_token;
    localStorage.setItem('access_token', data.access_token);
    localStorage.setItem('refresh_token', data.refresh_token);
    return true;
  }

  // Trigger Google OAuth Login popup window
  login(): Promise<AuthData> {
    return new Promise((resolve, reject) => {
      const width = 500;
      const height = 600;
      const left = window.screen.width / 2 - width / 2;
      const top = window.screen.height / 2 - height / 2;

      const popup = window.open(
        `${API_BASE}/api/auth/google`,
        'ISRO CIRT Authentication',
        `width=${width},height=${height},top=${top},left=${left}`
      );

      if (!popup) {
        reject(new Error('Popup blocker prevented OAuth window from opening.'));
        return;
      }

      let authCompleted = false;

      const messageListener = (event: MessageEvent) => {
        const allowedOrigins = [
          window.location.origin,
          API_BASE,
        ].filter(Boolean);
        if (!allowedOrigins.includes(event.origin)) return;
        if (event.data?.type === 'AUTH_SUCCESS') {
          authCompleted = true;
          clearInterval(timer);
          window.removeEventListener('message', messageListener);
          const data = event.data.data as AuthData;
          this.setTokens(data);
          resolve(data);
        }
      };

      window.addEventListener('message', messageListener);

      // Poll to check if popup closed without completing auth
      const timer = setInterval(() => {
        if (popup.closed && !authCompleted) {
          clearInterval(timer);
          window.removeEventListener('message', messageListener);
          reject(new Error('Login window closed before authentication was complete.'));
        }
      }, 500);
    });
  }

  async bypassLogin(username: string = 'divyanshxanshu'): Promise<AuthData> {
    const url = `${API_BASE}/api/auth/bypass?username=${encodeURIComponent(username)}`;
    const response = await fetch(url);
    if (!response.ok) {
      throw new Error(`Bypass authentication failed: ${response.statusText}`);
    }
    const data = await response.json() as AuthData;
    this.setTokens(data);
    return data;
  }


  async getCampaigns() {
    return this.request<{ campaigns: Array<{ name: string; slug: string; description: string; total_nodes: number; completed: number }> }>('/api/campaigns');
  }

  async getCampaignDetails(slug: string) {
    return this.request<{ campaign_slug: string; chapters: ChapterProgress[] }>(`/api/campaigns/${slug}`);
  }

  async getChapterDetails(slug: string, chapterNum: number) {
    return this.request<{ number: number; title: string; city: string; difficulty: string; commands: string[]; story_text: string }>(`/api/campaigns/${slug}/chapters/${chapterNum}`);
  }

  async submitFlag(slug: string, chapterNum: number, flag: string) {
    return this.request<{ correct: boolean; xp_earned: number; old_elo: number; new_elo: number; elo_diff: number; message: string }>(
      `/api/campaigns/${slug}/chapters/${chapterNum}/submit`,
      {
        method: 'POST',
        body: JSON.stringify({ flag }),
      }
    );
  }

  async startSandbox(chapterNumber: number) {
    return this.request<{ signed_url: string; expires_in: number }>('/api/sandbox/start', {
      method: 'POST',
      body: JSON.stringify({ chapter_number: chapterNumber }),
    });
  }

  async getUserProfile(username: string) {
    return this.request<UserProfile>(`/api/users/${username}/profile`);
  }

  async getLeaderboard() {
    return this.request<Array<{ username: string; elo: number; xp: number; rank: string }>>('/api/leaderboard');
  }

  async getWeeklyQuest() {
    return this.request<{ id: string; title: string; description: string; starts_at: string; ends_at: string; completed: boolean; attempts: number; hints_used: number }>('/api/competition/weekly');
  }

  async submitWeeklyQuest(flag: string) {
    return this.request<{ message: string; correct: boolean; score?: number }>('/api/competition/weekly/submit', {
      method: 'POST',
      body: JSON.stringify({ flag }),
    });
  }

  async getDailyChallenge() {
    return this.request<{ id: string; description: string; date: string; completed: boolean }>('/api/competition/daily');
  }

  async submitDailyChallenge(flag: string) {
    return this.request<{ message: string; correct: boolean; streak?: number }>('/api/competition/daily/submit', {
      method: 'POST',
      body: JSON.stringify({ flag }),
    });
  }
}

export const api = new ApiService();
