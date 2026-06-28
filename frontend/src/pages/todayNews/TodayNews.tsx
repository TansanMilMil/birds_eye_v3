import { Alert, Box, CircularProgress } from "@mui/material";
import { Masonry } from "@mui/lab";
import { MainArticle } from "../../share-components/article/MainArticle";
import { News } from "../../types/news";
import { useEffect, useState } from "react";
import { BirdsEyeApi } from "../../api/birdsEyeApi";
import { ThemeColorSetting } from "../../share-components/config/ThemeColorSetting";

export function TodayNews() {
  const [newsList, setNewsList] = useState<News[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [hasError, setHasError] = useState<boolean>(false);

  useEffect(() => {
    const fetchNews = async () => {
      setIsLoading(true);
      const today = new Date();

      for (let offset = 0; offset < 7; offset++) {
        const target = new Date(today);
        target.setDate(today.getDate() - offset);
        const targetDate = target.toISOString().slice(0, 10);

        try {
          const result = await BirdsEyeApi.getTodayNews(targetDate);
          if (result.news.length > 0) {
            const mapped: News[] = result.news.map((news) => {
              news.scrapedDateTime = new Date(
                Date.parse(news.scrapedDateTime)
              ).toLocaleString();
              return news;
            });
            setNewsList(mapped);
            setIsLoading(false);
            return;
          }
        } catch (err) {
          console.error(err);
          setIsLoading(false);
          setHasError(true);
          return;
        }
      }

      setIsLoading(false);
    };

    fetchNews();
  }, []);

  return (
    <div>
      <Box sx={{ maxWidth: "1200px", margin: "0 auto" }}>
        <Box sx={{ marginBottom: "2rem" }}>
          <ThemeColorSetting></ThemeColorSetting>
        </Box>
        {hasError && <Alert severity="error">network error...</Alert>}
        {isLoading && (
          <Box sx={{ textAlign: "center", margin: "3rem" }}>
            <CircularProgress color="secondary" />
          </Box>
        )}
        {!isLoading && (
          <Masonry
            columns={{ xs: 1, sm: 2, md: 3 }}
            spacing={{ xs: 1, sm: 1, md: 1 }}
          >
            {newsList.map((news, i) => (
              <MainArticle
                key={i}
                news={news}
                isDisplayReactions={true}
              ></MainArticle>
            ))}
          </Masonry>
        )}
      </Box>
    </div>
  );
}
