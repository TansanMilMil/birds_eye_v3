import { Box, IconButton } from "@mui/material";
import { useEffect } from "react";
import { useDispatch } from "react-redux";
import { changeThemePalette } from "../../app/themeSlice";
import { LocalStorageKey } from "../../localStorage/localStorageKey";
import { MuiTheme } from "../../mui-theme";
import { MuiThemeNames } from "../../mui-theme-names";
import ForestIcon from "@mui/icons-material/Forest";
import BeachAccessIcon from "@mui/icons-material/BeachAccess";
import MosqueIcon from "@mui/icons-material/Mosque";
import LightModeIcon from "@mui/icons-material/LightMode";
import NightsStayIcon from "@mui/icons-material/NightsStay";
import LocalFloristIcon from "@mui/icons-material/LocalFlorist";
import WavesIcon from "@mui/icons-material/Waves";
import WbTwilightIcon from "@mui/icons-material/WbTwilight";
import ParkIcon from "@mui/icons-material/Park";

export function ThemeColorSetting() {
  const dispatch = useDispatch();

  useEffect(() => {
    const themeName: string | null = localStorage.getItem(
      LocalStorageKey.ThemeName
    );
    if (themeName) {
      MuiTheme.AssertTheme(themeName);
      dispatch(changeThemePalette(themeName));
    }
  }, [dispatch]);

  const changeTheme = (themeName: string) => {
    MuiTheme.AssertTheme(themeName);
    localStorage.setItem(LocalStorageKey.ThemeName, themeName);
    dispatch(changeThemePalette(themeName));
  };

  return (
    <div>
      <Box sx={{ textAlign: "center" }}>
        <IconButton onClick={() => changeTheme(MuiThemeNames.Nature)}>
          <ForestIcon />
        </IconButton>
        <IconButton onClick={() => changeTheme(MuiThemeNames.Aqua)}>
          <BeachAccessIcon />
        </IconButton>
        <IconButton onClick={() => changeTheme(MuiThemeNames.Horror)}>
          <MosqueIcon />
        </IconButton>
        <IconButton onClick={() => changeTheme(MuiThemeNames.Sunny)}>
          <LightModeIcon />
        </IconButton>
        <IconButton onClick={() => changeTheme(MuiThemeNames.Night)}>
          <NightsStayIcon />
        </IconButton>
        <IconButton onClick={() => changeTheme(MuiThemeNames.Sakura)}>
          <LocalFloristIcon />
        </IconButton>
        <IconButton onClick={() => changeTheme(MuiThemeNames.Ocean)}>
          <WavesIcon />
        </IconButton>
        <IconButton onClick={() => changeTheme(MuiThemeNames.Sunset)}>
          <WbTwilightIcon />
        </IconButton>
        <IconButton onClick={() => changeTheme(MuiThemeNames.Forest)}>
          <ParkIcon />
        </IconButton>
      </Box>
    </div>
  );
}
