<script lang="ts">
  /**
   * This file bootstraps the RouterContext with the defined routes.
   * Ensure to wrap the application with <Routing> in App.svelte.
   *
   * Then use <RouterView /> to render the matched route component in
   * the desired location.
   */
  import type { RouterOptions } from "@dvcol/svelte-simple-router/models";

  import { RouterContext } from "@dvcol/svelte-simple-router/components";
  import type { Snippet } from "svelte";

  function withBaseURL(path: string) {
    return (import.meta.env.BASE_URL || "/") + path;
  }

  interface Props {
    children: Snippet;
  }

  const { children }: Props = $props();

  const options: RouterOptions = {
    routes: [
      {
        name: "home",
        path: withBaseURL(""),
        component: () => import("./pages/Home.svelte"),
      },
      {
        name: "images",
        path: withBaseURL("images"),
        component: () => import("./pages/Images.svelte"),
      },
      {
        name: "devices",
        path: withBaseURL("devices"),
        component: () => import("#/pages/Devices/Devices.svelte"),
      },
      {
        name: "sources",
        path: withBaseURL("sources"),
        component: () => import("#/pages/Sources/Sources.svelte"),
      },
      {
        name: "jobs",
        path: withBaseURL("jobs"),
        component: () => import("./pages/Jobs.svelte"),
      },
      {
        // TODO: create a NotFound.svelte page
        name: "not-found",
        path: withBaseURL("*"),
        redirect: {
          name: "home",
        },
      },
    ],
  };
</script>

<RouterContext {options}>
  {@render children()}
</RouterContext>
