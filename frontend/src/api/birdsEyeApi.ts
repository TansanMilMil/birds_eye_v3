import { Env } from "../env";
import {
  GetReactionsResponse,
  GetTodayNewsResponse,
  GetTrendsResponse,
} from "../types/birdsEyeApi";

export class BirdsEyeApi {
  private static readonly API_ENDPOINT: string = Env.BirdsEyeApiEndpoint();

  private static async get<T>(path: string): Promise<T> {
    const res = await fetch(`${this.API_ENDPOINT}${path}`);
    if (!res.ok) {
      throw new Error(`GET ${path} failed: ${res.status}`);
    }
    return res.json() as Promise<T>;
  }

  public static async getTodayNews(targetDateTime: string): Promise<GetTodayNewsResponse> {
    return this.get(`/news/${encodeURIComponent(targetDateTime)}`);
  }

  public static async getTrends(): Promise<GetTrendsResponse> {
    return this.get("/news/trends");
  }

  public static async getReactions(id: number): Promise<GetReactionsResponse> {
    return this.get(`/news/news-reactions/${id}`);
  }
}
