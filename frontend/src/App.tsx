import {
  BottomNavigation,
  BottomNavigationAction,
  Box,
  Theme,
} from "@mui/material";
import { TodayNews } from "./pages/todayNews/TodayNews";
import title from "./images/logo.png";
import { useEffect, useState } from "react";
import NewspaperIcon from "@mui/icons-material/Newspaper";
import BubbleChartIcon from "@mui/icons-material/BubbleChart";
import { Trends } from "./pages/trends/Trends";
import { Route, Routes } from "react-router-dom";
import { useNavigate } from "react-router-dom";
import styles from "./App.module.css";
import { useSelector } from "react-redux";
import { createTheme, ThemeProvider } from "@mui/material/styles";
import { MuiTheme } from "./mui-theme";
import { RootState } from "./app/store";

export function App() {
  const navigate = useNavigate();
  const [navigationValue, setNavigationValue] = useState<string>();
  const theme: Theme = useSelector((state: RootState) => {
    return createTheme({
      palette: MuiTheme.GetTheme(state.theme.themeName),
    });
  });

  useEffect(() => {
    if (window.location.pathname === "/") {
      navigate("/news");
    }
    setNavigationValue(window.location.pathname);
  }, [navigate]);

  return (
    <ThemeProvider theme={theme}>
      <Box sx={{ backgroundColor: "primary.main" }}>
        <div className={styles.appBase}>
          <h1 className={styles.title}>
            <img src={title} alt="title" />
          </h1>

          <Routes>
            <Route path="/news" element={<TodayNews></TodayNews>} />
            <Route path="/trends" element={<Trends></Trends>} />
          </Routes>

          <div className={styles.credit}>
            <div>
              ロゴは{" "}
              <a
                href="https://www.designevo.com/jp/"
                title="無料オンラインロゴメーカー"
              >
                DesignEvo
              </a>{" "}
              ロゴメーカーさんに作られる
            </div>
          </div>

          <BottomNavigation
            showLabels
            value={navigationValue}
            onChange={(event, newValue) => {
              setNavigationValue(newValue);
              navigate(newValue);
            }}
            sx={{
              position: "fixed",
              bottom: "0",
              left: "0",
              width: "100%",
              zIndex: 100,
            }}
          >
            <BottomNavigationAction
              label="News"
              icon={<NewspaperIcon />}
              value="/news"
            />
            <BottomNavigationAction
              label="Trends"
              icon={<BubbleChartIcon />}
              value="/trends"
            />
          </BottomNavigation>
        </div>
      </Box>
    </ThemeProvider>
  );
}
