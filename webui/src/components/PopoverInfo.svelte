<script lang="ts">
  import { Popover, type WithoutChildren } from "bits-ui";
  import type { Snippet } from "svelte";
  import type { ClassValue } from "svelte/elements";
  import { theme } from "#/store/theme";
  import IconX from "@lucide/svelte/icons/x";
  import IconInfo from "@lucide/svelte/icons/info";

  type Props = Popover.RootProps & {
    trigger?: Snippet;
    title?: string | Snippet;
    class?: ClassValue;
    contentProps?: WithoutChildren<Popover.ContentProps>;
  } & (
      | { children: Snippet; markdown?: never }
      | { children?: never; markdown: string }
    );
  let {
    open = $bindable(false),
    children,
    trigger,
    title,
    markdown,
    class: className,
    contentProps = {},
    ...restProps
  }: Props = $props();
</script>

<Popover.Root bind:open {...restProps}>
  <Popover.Trigger>
    {#if trigger}
      {@render trigger()}
    {:else}
      <button type="button" class="btn btn-circle btn-ghost btn-xs">
        <IconInfo />
      </button>
    {/if}
  </Popover.Trigger>
  <Popover.Portal>
    <Popover.Content
      {...contentProps}
      class={"z-[9999] bg-transparent max-w-[90vw] max-h-[60vh]" +
        (className ? " " + className : "")}
      data-theme={$theme}
    >
      <Popover.Arrow class="text-primary" />
      <div class="card bg-base-300 border border-base-200 max-h-[50vh]">
        <div class="card-body">
          <div class="card-title justify-between items-start">
            {#if title}
              {#if typeof title === "string"}
                <span>{title}</span>
              {:else}
                {@render title()}
              {/if}
            {:else}
              <span class="font-bold">Info</span>
            {/if}
            <Popover.Close class="btn btn-square btn-ghost btn-xs">
              <IconX />
            </Popover.Close>
          </div>
          <div
            class="p-4 border border-base-100 max-h-[35vh] overflow-auto"
            style="box-shadow: inset 0 6px 12px rgba(0, 0, 0, 0.15), inset 0 2px 4px rgba(0, 0, 0, 0.1);"
          >
            {#if children}
              {@render children()}
            {:else if markdown && open}
              {#await import ("#/components/MarkdownText.svelte") then { default: MarkdownText }}
                <MarkdownText text={markdown} />
              {/await}
            {/if}
          </div>
        </div>
      </div>
    </Popover.Content>
  </Popover.Portal>
</Popover.Root>
