import { Alert, Box, CircularProgress, IconButton, Typography } from "@mui/material";
import { Masonry } from "@mui/lab";
import ArrowBackIosNewIcon from "@mui/icons-material/ArrowBackIosNew";
import ArrowForwardIosIcon from "@mui/icons-material/ArrowForwardIos";
import { MainArticle } from "../../share-components/article/MainArticle";
import { News } from "../../types/news";
import { useEffect, useState } from "react";
import { BirdsEyeApi } from "../../api/birdsEyeApi";
import { ThemeColorSetting } from "../../share-components/config/ThemeColorSetting";

const toDateKey = (date: Date): string => date.toISOString().slice(0, 10);

const getToday = (): Date => {
  const today = new Date();
  today.setHours(0, 0, 0, 0);
  return today;
};

export function TodayNews() {
  const [newsList, setNewsList] = useState<News[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [hasError, setHasError] = useState<boolean>(false);
  const [targetDate, setTargetDate] = useState<Date | null>(null);

  const fetchNewsForDate = async (date: Date) => {
    setIsLoading(true);
    setHasError(false);
    try {
      const result = await BirdsEyeApi.getTodayNews(toDateKey(date));
      const mapped: News[] = result.news.map((news) => {
        news.scrapedDateTime = new Date(
          Date.parse(news.scrapedDateTime)
        ).toLocaleString();
        return news;
      });
      setNewsList(mapped);
    } catch (err) {
      console.error(err);
      setHasError(true);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    const findLatestAvailableDate = async () => {
      setIsLoading(true);
      const today = getToday();

      for (let offset = 0; offset < 7; offset++) {
        const target = new Date(today);
        target.setDate(today.getDate() - offset);

        try {
          const result = await BirdsEyeApi.getTodayNews(toDateKey(target));
          if (result.news.length > 0) {
            const mapped: News[] = result.news.map((news) => {
              news.scrapedDateTime = new Date(
                Date.parse(news.scrapedDateTime)
              ).toLocaleString();
              return news;
            });
            setNewsList(mapped);
            setTargetDate(target);
            setIsLoading(false);
            return;
          }
        } catch (err) {
          console.error(err);
          setHasError(true);
          setIsLoading(false);
          return;
        }
      }

      setTargetDate(today);
      setIsLoading(false);
    };

    findLatestAvailableDate();
  }, []);

  const handlePrevDay = () => {
    if (!targetDate) return;
    const prev = new Date(targetDate);
    prev.setDate(targetDate.getDate() - 1);
    setTargetDate(prev);
    fetchNewsForDate(prev);
  };

  const handleNextDay = () => {
    if (!targetDate) return;
    const next = new Date(targetDate);
    next.setDate(targetDate.getDate() + 1);
    setTargetDate(next);
    fetchNewsForDate(next);
  };

  const isNextDisabled = targetDate
    ? toDateKey(targetDate) >= toDateKey(getToday())
    : true;

  return (
    <div>
      <Box sx={{ maxWidth: "1200px", margin: "0 auto" }}>
        <Box sx={{ marginBottom: "2rem" }}>
          <ThemeColorSetting></ThemeColorSetting>
        </Box>
        {targetDate && (
          <Box
            sx={{
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              gap: 1,
              marginBottom: "1rem",
            }}
          >
            <IconButton onClick={handlePrevDay} disabled={isLoading}>
              <ArrowBackIosNewIcon />
            </IconButton>
            <Typography variant="h6">
              {targetDate.toLocaleDateString("ja-JP", {
                year: "numeric",
                month: "long",
                day: "numeric",
              })}
            </Typography>
            <IconButton
              onClick={handleNextDay}
              disabled={isLoading || isNextDisabled}
            >
              <ArrowForwardIosIcon />
            </IconButton>
          </Box>
        )}
        {hasError && <Alert severity="error">network error...</Alert>}
        {isLoading && (
          <Box sx={{ textAlign: "center", margin: "3rem" }}>
            <CircularProgress color="secondary" />
          </Box>
        )}
        {!isLoading && !hasError && newsList.length === 0 && (
          <Typography sx={{ textAlign: "center", margin: "3rem" }}>
            この日のニュースはありません
          </Typography>
        )}
        {!isLoading && newsList.length > 0 && (
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
