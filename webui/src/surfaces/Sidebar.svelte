<script lang="ts">
  import { Sidebar, SidebarGroup, SidebarItem } from "flowbite-svelte";
  import { useRouter } from "@dvcol/svelte-simple-router/router";
  import IconHouse from "@lucide/svelte/icons/house";
  import IconImage from "@lucide/svelte/icons/image";
  import IconMonitorSmartphone from "@lucide/svelte/icons/monitor-smartphone";
  import IconDatabase from "@lucide/svelte/icons/database";
  import IconTask from "@lucide/svelte/icons/square-check-big";
  import IconChevronsDown from "@lucide/svelte/icons/chevrons-down";
  import IconSettings from "@lucide/svelte/icons/settings";

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

<Sidebar
  backdrop={isMobile}
  {isOpen}
  closeSidebar={() => (isOpen = false)}
  class="w-24 static"
  divClass="h-screen flex flex-col justify-between w-full px-0"
>
  <SidebarGroup>
    <SidebarItem
      aClass="flex-col justify-center items-center gap-2"
      spanClass="mx-0"
      label="Claw"
    >
      {#snippet icon()}
        <IconChevronsDown />
      {/snippet}
    </SidebarItem>
  </SidebarGroup>
  <SidebarGroup>
    {#each items as item (item.id)}
      <SidebarItem
        label={item.label}
        class="hover:cursor-pointer"
        aClass="flex-col justify-center items-center gap-2"
        onclick={() => router.push({ name: item.id })}
        spanClass="mx-0"
        active={location?.name === item.id}
      >
        {#snippet icon()}
          <item.Icon />
        {/snippet}
      </SidebarItem>
    {/each}
  </SidebarGroup>
  <SidebarGroup>
    <SidebarItem
      aClass="flex-col justify-center items-center gap-2"
      spanClass="mx-0"
      label="Settings"
    >
      {#snippet icon()}
        <IconSettings />
      {/snippet}
    </SidebarItem>
  </SidebarGroup>
</Sidebar>
