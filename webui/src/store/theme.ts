import { writable, derived } from "svelte/store";

export const themes = [
  { value: "nord", name: "Nord" },
  { value: "dracula", name: "Dracula" },
  { value: "light", name: "Light" },
  { value: "dark", name: "Dark" },
] as const;

export type ThemeName = (typeof themes)[number]["value"];

const current = writable(
  (localStorage.getItem("theme") || themes[0].value) as ThemeName,
);

/**
 * A Svelte store that holds the current theme.
 *
 * If invalid theme is found in localStorage or set somehow set incorrectly by user,
 * it defaults to the first theme in the themes array.
 */
export const theme = derived(current, (value) => {
  if (!themes.find((v) => v.value === value)) {
    return themes[0];
  }
  return value;
});

/**
 * Sets the current theme (and updates the theme variable) and saves it to localStorage.
 *
 * On website reload (e.g. by refresh), the theme is read from localStorage, keeping the theme persistent.
 *
 * Setting theme unavailable in the themes array will throw an error.
 */
export function setTheme(newTheme: ThemeName) {
  if (!themes.find((v) => v.value === newTheme)) {
    throw new Error(`Unknown theme: ${newTheme}`);
  }
  current.set(newTheme);
  localStorage.setItem("theme", newTheme);
}
