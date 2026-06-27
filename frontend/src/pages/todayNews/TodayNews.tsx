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
    if (newsList.length === 0) {
      setIsLoading(true);
      const today = new Date();
      const year = today.getFullYear();
      const month = String(today.getMonth() + 1).padStart(2, "0");
      const day = String(today.getDate()).padStart(2, "0");
      const targetDate = `${year}-${month}-${day}`;

      BirdsEyeApi.getTodayNews(targetDate)
        .then((result) => {
          const newsList: News[] = result.news.map((news) => {
            news.scrapedDateTime = new Date(
              Date.parse(news.scrapedDateTime)
            ).toLocaleString();
            return news;
          });
          setNewsList(newsList);
          setIsLoading(false);
        })
        .catch((err) => {
          console.error(err);
          setIsLoading(false);
          setHasError(true);
        });
    }
  }, [newsList.length]);

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
