import { PaletteOptions } from "@mui/material";
import { amber, blue, blueGrey, brown, cyan, deepOrange, green, indigo, lightBlue, orange, pink, purple, red, teal } from "@mui/material/colors";
import { MuiThemeNames } from "./mui-theme-names";

export class MuiTheme {
  public static readonly themeItemList = [
    {
      themeName: MuiThemeNames.Nature,
      themePalette: {
          primary: teal,
          secondary: pink,
      },
    },
    {
      themeName: MuiThemeNames.Aqua,
      themePalette: {
          primary: blue,
          secondary: deepOrange,
      },
    },
    {
      themeName: MuiThemeNames.Horror,
      themePalette: {
          primary: purple,
          secondary: red,
      },
    },
    {
      themeName: MuiThemeNames.Sunny,
      themePalette: {
          primary: orange,
          secondary: blue,
      },
    },
    {
      themeName: MuiThemeNames.Night,
      themePalette: {
          primary: indigo,
          secondary: blueGrey,
      },
    },
    {
      themeName: MuiThemeNames.Sakura,
      themePalette: {
          primary: pink,
          secondary: purple,
      },
    },
    {
      themeName: MuiThemeNames.Ocean,
      themePalette: {
          primary: lightBlue,
          secondary: cyan,
      },
    },
    {
      themeName: MuiThemeNames.Sunset,
      themePalette: {
          primary: deepOrange,
          secondary: amber,
      },
    },
    {
      themeName: MuiThemeNames.Forest,
      themePalette: {
          primary: green,
          secondary: brown,
      },
    },
  ];

  public static readonly DefaultThemeName = MuiThemeNames.Nature;

  public static GetDefaultTheme(): PaletteOptions {
    return this.themeItemList.find(x => x.themeName === this.DefaultThemeName)!.themePalette;
  }

  public static GetTheme(themeName: string): PaletteOptions {
    this.AssertTheme(themeName);
    const themeItem = this.themeItemList.find(x => x.themeName === themeName);
    return themeItem!.themePalette;
  }

  public static AssertTheme(themeName: string | null): void {
    if (!themeName || !this.themeItemList.some(x => x.themeName === themeName)) {
      throw new Error(`themeName: ${themeName} is not exist.`);
    }
  }
}