import {
  Alert,
  Box,
  CircularProgress,
  Divider,
  IconButton,
  Typography,
} from "@mui/material";
import { Masonry } from "@mui/lab";
import ArrowBackIosNewIcon from "@mui/icons-material/ArrowBackIosNew";
import ArrowForwardIosIcon from "@mui/icons-material/ArrowForwardIos";
import FiberManualRecordIcon from "@mui/icons-material/FiberManualRecord";
import LocalFireDepartmentIcon from "@mui/icons-material/LocalFireDepartment";
import { MainArticle } from "../../share-components/article/MainArticle";
import { News } from "../../types/news";
import { useEffect, useState } from "react";
import { BirdsEyeApi } from "../../api/birdsEyeApi";
import { ThemeColorSetting } from "../../share-components/config/ThemeColorSetting";

const PICKUP_MIN_IMPORTANCE = 0.6;
const PICKUP_LIMIT = 3;

const toDateKey = (date: Date): string => {
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const day = String(date.getDate()).padStart(2, "0");
  return `${date.getFullYear()}-${month}-${day}`;
};

const getToday = (): Date => {
  const today = new Date();
  today.setHours(0, 0, 0, 0);
  return today;
};

const shiftDate = (date: Date, days: number): Date => {
  const shifted = new Date(date);
  shifted.setDate(date.getDate() + days);
  return shifted;
};

const fetchNews = async (date: Date): Promise<News[]> => {
  const result = await BirdsEyeApi.getTodayNews(toDateKey(date));
  return result.news.map((news) => {
    news.scrapedDateTime = new Date(
      Date.parse(news.scrapedDateTime)
    ).toLocaleString();
    return news;
  });
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
      setNewsList(await fetchNews(date));
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
      try {
        for (let offset = 0; offset < 7; offset++) {
          const target = shiftDate(today, -offset);
          const news = await fetchNews(target);
          if (news.length > 0) {
            setNewsList(news);
            setTargetDate(target);
            return;
          }
        }
        setTargetDate(today);
      } catch (err) {
        console.error(err);
        setHasError(true);
      } finally {
        setIsLoading(false);
      }
    };

    findLatestAvailableDate();
  }, []);

  const moveDay = (days: number) => {
    if (!targetDate) return;
    const next = shiftDate(targetDate, days);
    setTargetDate(next);
    fetchNewsForDate(next);
  };

  const pickupNews = newsList
    .filter((news) => news.importance >= PICKUP_MIN_IMPORTANCE)
    .slice(0, PICKUP_LIMIT);
  const otherNews = newsList.filter((news) => !pickupNews.includes(news));

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
            <IconButton onClick={() => moveDay(-1)} disabled={isLoading}>
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
              onClick={() => moveDay(1)}
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
        {!isLoading && pickupNews.length > 0 && (
          <Box sx={{ marginBottom: "2rem" }}>
            <Box
              sx={{
                display: "flex",
                alignItems: "center",
                gap: 1,
                marginBottom: "0.75rem",
                color: "common.white",
              }}
            >
              <LocalFireDepartmentIcon />
              <Typography
                variant="h6"
                component="h2"
                sx={{ fontSize: "1rem", fontWeight: "bold", letterSpacing: "0.1em" }}
              >
                ピックアップ
              </Typography>
            </Box>
            <Masonry
              columns={{ xs: 1, sm: 2, md: 3 }}
              spacing={{ xs: 1, sm: 1, md: 1 }}
            >
              {pickupNews.map((news) => (
                <MainArticle
                  key={news.id}
                  news={news}
                  isDisplayReactions={true}
                ></MainArticle>
              ))}
            </Masonry>
          </Box>
        )}
        {!isLoading && pickupNews.length > 0 && otherNews.length > 0 && (
          <Divider
            textAlign="center"
            sx={{
              marginBottom: "1.5rem",
              color: "common.white",
              "&::before, &::after": { borderColor: "common.white" },
            }}
          >
            <FiberManualRecordIcon sx={{ fontSize: "0.7rem" }} />
          </Divider>
        )}
        {!isLoading && otherNews.length > 0 && (
          <Masonry
            columns={{ xs: 1, sm: 2, md: 3 }}
            spacing={{ xs: 1, sm: 1, md: 1 }}
          >
            {otherNews.map((news) => (
              <MainArticle
                key={news.id}
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
