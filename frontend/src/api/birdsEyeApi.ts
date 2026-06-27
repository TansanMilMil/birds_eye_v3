import axios from "axios";
import { Env } from "../env";
import {
  GetReactionsResponse,
  GetTodayNewsResponse,
  GetTrendsResponse,
} from "../types/birdsEyeApi";

export class BirdsEyeApi {
  private static readonly API_ENDPOINT: string = Env.BirdsEyeApiEndpoint();

  public static async getTodayNews(): Promise<GetTodayNewsResponse> {
    const res = await axios.get(`${this.API_ENDPOINT}/news/today-news`);
    return res.data;
  }

  public static async getTrends(): Promise<GetTrendsResponse> {
    const res = await axios.get(`${this.API_ENDPOINT}/news/trends`);
    return res.data;
  }

  public static async getReactions(id: number): Promise<GetReactionsResponse> {
    const res = await axios.get(
      `${this.API_ENDPOINT}/news/news-reactions/${id}`
    );
    return res.data;
  }
}
