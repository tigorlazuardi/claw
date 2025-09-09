<script lang="ts">
  import { createPopover, melt } from "@melt-ui/svelte";
  import { fade } from "svelte/transition";
  import IconX from "@lucide/svelte/icons/x";
  import IconInfo from "@lucide/svelte/icons/info";

  interface Props {
    /** The markdown content to display inside the popover. */
    content: string;
    headerText?: string;
  }

  const { content: md, headerText = "Info" }: Props = $props();

  const {
    elements: { content, arrow, close, trigger },
    states: { open },
  } = createPopover({
    forceVisible: true,
  });
</script>

<button
  class="btn btn-square btn-ghost btn-xs"
  use:melt={$trigger}
  aria-label="Info"
>
  <IconInfo />
</button>

{#if $open}
  <div
    data-theme="nord"
    class="card z-[9999]"
    use:melt={$content}
    transition:fade={{ duration: 100 }}
  >
    <div use:melt={$arrow}></div>
    <div class="card-title justify-between px-4 pt-2">
      <h3 class="text-3xl">{headerText}</h3>
      <button class="btn btn-square btn-ghost" use:melt={$close}>
        <IconX />
      </button>
    </div>
    <div class="divider"></div>

    {#await import("./MarkdownText.svelte")}
      <div class="h-full w-full flex justify-center items-center">
        <span class="loading loading-ring text-primary w-[6rem] h-auto"></span>
      </div>
    {:then { default: MarkdownText }}
      <MarkdownText text={md} />
    {/await}
  </div>
{/if}
