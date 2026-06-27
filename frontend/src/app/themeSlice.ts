import {
  createSlice,
  PayloadAction,
  SliceCaseReducers,
  SliceSelectors,
} from "@reduxjs/toolkit";
import { MuiTheme } from "../mui-theme";

export type ThemeState = {
  themeName: string;
};

export const themeSlice = createSlice<
  ThemeState,
  SliceCaseReducers<ThemeState>,
  string,
  SliceSelectors<ThemeState>,
  string
>({
  name: "theme",
  initialState: {
    themeName: MuiTheme.DefaultThemeName,
  },
  reducers: {
    changeThemePalette: (state, action: PayloadAction<string>) => {
      return { ...state, themeName: action.payload };
    },
  },
});

export const { changeThemePalette } = themeSlice.actions;

export default themeSlice.reducer;
