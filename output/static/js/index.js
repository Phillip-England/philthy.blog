// function getCookie(name) {
//   const value = `; ${document.cookie}`;
//   const parts = value.split(`; ${name}=`);
//   if (parts.length === 2) return parts.pop().split(";").shift();
// }

// function setCookie(name, value, days) {
//   const date = new Date();
//   date.setTime(date.getTime() + days * 24 * 60 * 60 * 1000);
//   const expires = `expires=${date.toUTCString()}`;
//   document.cookie = `${name}=${value}; ${expires}; path=/`;
// }

// document.addEventListener("DOMContentLoaded", function () {
//   let theme = getCookie("theme");
//   if (!theme) {
//     const prefersDarkScheme = window.matchMedia(
//       "(prefers-color-scheme: dark)",
//     ).matches;
//     theme = prefersDarkScheme ? "dark" : "light";
//     setCookie("theme", theme, 365);
//   }
//   if (theme === "dark") {
//     document.documentElement.classList.add("dark");
//   }
// });
