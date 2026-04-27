export type Role = "admin" | "manager" | "user";

export interface Credentials {
  email: string;
  password: string;
}

export interface AuthResponse {
  token: string;
  user: UserProfile;
  message?: string;
}

export interface UserProfile {
  id: number;
  email: string;
  role: string;
  status: number;
  createdAt?: string;
}
