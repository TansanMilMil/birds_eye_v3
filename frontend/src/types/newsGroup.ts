import { News } from "./news";

export type NewsGroup = {
    sourceBy: string;
    scrapedUrl: string;
    news: News[];
}