export type ReactionSentiment = {
    positive: number;
    neutral: number;
    negative: number;
}

export type News = {
    id: number;
    title: string;
    description: string;
    summarizedText: string;
    sourceBy: string;
    category: string;
    importance: number;
    scrapedUrl: string;
    scrapedDateTime: string;
    articleUrl: string;
    articleImageUrl: string;
    reactionCount: number;
    reactionSentiment: ReactionSentiment;
}
