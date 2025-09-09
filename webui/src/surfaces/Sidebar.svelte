<script lang="ts">
  import {
    Sidebar,
    SidebarBrand,
    SidebarButton,
    SidebarCta,
    SidebarDropdownItem,
    SidebarDropdownWrapper,
    SidebarGroup,
    SidebarItem,
    SidebarWrapper,
  } from "flowbite-svelte";
  import {
    HomeOutline,
    ImageOutline,
    DesktopPcOutline,
    DatabaseOutline,
    CheckOutline,
    CogOutline,
    BarsOutline,
  } from "flowbite-svelte-icons";
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

  let innerWidth = $state(window.innerWidth);
  let isMobile = $derived(innerWidth < 640);
  let isOpen = $state(false);
  let isExpanded = $state(false);
  $effect(() => {
    // Use $effect so only run once on component mount.
    if (window.innerWidth >= 640) {
      isOpen = true;
    }

    if (window.innerWidth >= 1280) {
      isExpanded = true;
    }
  });

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

<svelte:window bind:innerWidth />

{#if isMobile}
  {#await import("flowbite-svelte") then { SidebarButton }}
    <SidebarButton onclick={() => (isOpen = !isOpen)} class="mb-2" />
  {/await}
{/if}
