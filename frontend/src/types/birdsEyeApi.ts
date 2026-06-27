import { News } from "./news";
import { NewsReaction } from "./newsReaction";

export interface GetTodayNewsResponse {
  news: News[];
}

export interface GetTrendsResponse {
  trends: News[];
}

export interface GetReactionsResponse {
  reactions: NewsReaction[];
}
