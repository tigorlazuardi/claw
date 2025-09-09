<script lang="ts">
  import { useRouter } from "@dvcol/svelte-simple-router/router";
  import { link } from "@dvcol/svelte-simple-router";
  import IconHouse from "@lucide/svelte/icons/house";
  import IconImage from "@lucide/svelte/icons/image";
  import IconMonitorSmartphone from "@lucide/svelte/icons/monitor-smartphone";
  import IconDatabase from "@lucide/svelte/icons/database";
  import IconTask from "@lucide/svelte/icons/layout-list";
  import IconSettings from "@lucide/svelte/icons/settings";
  import { type Snippet } from "svelte";

  interface Props {
    children: Snippet;
  }

  let { children }: Props = $props();

  function withBaseURL(str: string) {
    return (import.meta.env.BASE_URL || "/") + str;
  }

  const router = useRouter();
  const { location } = $derived(router);

  const items = [
    {
      id: "home",
      label: "Home",
      Icon: IconHouse,
      href: withBaseURL(""),
    },
    {
      id: "images",
      label: "Images",
      Icon: IconImage,
      href: withBaseURL("images"),
    },
    {
      id: "devices",
      label: "Devices",
      Icon: IconMonitorSmartphone,
      href: withBaseURL("devices"),
    },
    {
      id: "sources",
      label: "Sources",
      Icon: IconDatabase,
      href: withBaseURL("sources"),
    },
    {
      id: "jobs",
      label: "Jobs",
      Icon: IconTask,
      href: withBaseURL("jobs"),
    },
  ];
</script>

<div class="drawer drawer-open">
  <input id="sidebar-toggle" type="checkbox" class="drawer-toggle" />
  <div class="drawer-content">
    {@render children()}
  </div>
  <div class="drawer-side">
    <label
      for="my-drawer"
      aria-label="close sidebar"
      class="drawer-overlay"
    ></label>
    <aside class="flex flex-col justify-between h-screen bg-primary">
      <!-- Spacer to push nav-links to center -->
      <div></div>

      <div id="nav-links" class="m-auto flex flex-col">
        {#each items as item (item.id)}
          <a
            href={item.href}
            class={{
              "btn btn-sm btn-primary flex-col h-[4rem] rounded-none": true,
              "btn-secondary": location?.name === item.id,
            }}
            use:link
          >
            <item.Icon />
            {item.label}
          </a>
        {/each}
      </div>
      <button
        class="btn btn-sm mx-auto w-full btn-primary h-[4rem] rounded-none"
      >
        <IconSettings />
      </button>
    </aside>
  </div>
</div>
