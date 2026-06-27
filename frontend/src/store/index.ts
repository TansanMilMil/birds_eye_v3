import { green, lightGreen, orange, pink, teal } from '@mui/material/colors';
import { createTheme, PaletteOptions, Theme } from '@mui/material';
import { createStore } from "redux";

const initState: { theme: Theme } = {
    theme: createTheme({
        palette: {
            primary: teal,
            secondary: pink,
        }
    })
};

const reducer = (state = initState) => {
    return state;
};

const store = createStore(reducer);

export default store;