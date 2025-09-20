import { writable } from "svelte/store";

export let isMobile = writable(window.innerWidth <= 640);
