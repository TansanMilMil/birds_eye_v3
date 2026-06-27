import { TrendData } from "./trendData";

export type TrendSeries = {
    type: string;
    data: TrendData[];
    name: string;
}