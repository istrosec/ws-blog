import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import reportWebVitals from './reportWebVitals';
import {createTheme, CssBaseline, ThemeProvider} from "@mui/material";
import istrosecSecondary from '@mui/material/colors/amber';

export const istrosec = {
    50: '#e1faf7',
    100: '#bdfff7',
    200: '#81ebde',
    300: '#5cb5ab',
    400: '#57ffeb',
    500: '#26ffe6',
    600: '#22e6cf',
    700: '#1dbfac',
    800: '#1bb5a3',
    900: '#138075',
    A100: '#85fff1',
    A200: '#4affea',
    A400: '#26ffe6',
    A700: '#1fccb8',
    main: '#26ffe6'
};

const mainTheme = createTheme({
    palette: {
        mode: 'dark',
        primary: istrosec,
        secondary: istrosecSecondary
    },
});

ReactDOM.render(
    <ThemeProvider theme={mainTheme}>
        <CssBaseline/>
        <React.StrictMode>
            <App/>
        </React.StrictMode>
    </ThemeProvider>,
    document.getElementById('root')
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
