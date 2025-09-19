<script lang="ts">
  import type { Snippet } from "svelte";
  import { Dialog, type WithoutChild } from "bits-ui";
  import IconX from "@lucide/svelte/icons/x";
  import { theme } from "#/store/theme";
  import { isMobile } from "#/store/isMobile";
  import type { ClassValue } from "svelte/elements";

  type Props = Dialog.RootProps & {
    trigger?: Snippet;
    title: Snippet;
    description: Snippet;
    actions?: Snippet;
    contentProps?: WithoutChild<Dialog.ContentProps>;
    class?: ClassValue;
    // ...other component props if you wish to pass them
  };

  let {
    open = $bindable(false),
    children,
    trigger,
    contentProps = {},
    title,
    description,
    actions,
    class: className,
    ...restProps
  }: Props = $props();
</script>

<Dialog.Root bind:open {...restProps}>
  {#if trigger}
    <Dialog.Trigger>
      {@render trigger()}
    </Dialog.Trigger>
  {/if}
  <Dialog.Portal to={isMobile ? document.body : "main"}>
    <Dialog.Overlay class="fixed h-screen w-screen bg-black/40" />
    <Dialog.Content
      {...contentProps}
      class={"z-[9999] fixed top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 bg-base-200 p-6 rounded-lg shadow-lg" +
        (className ? " " + className : "")}
      data-theme={$theme}
    >
      <Dialog.Title class="flex justify-between items-center mb-4">
        {@render title()}
        <Dialog.Close class="btn btn-square btn-ghost btn-xs">
          <IconX />
        </Dialog.Close>
      </Dialog.Title>
      <div class="divider"></div>
      <Dialog.Description>
        {@render description()}
      </Dialog.Description>
      {@render children?.()}
      {#if actions}
        <div class="divider"></div>
        <div class="flex justify-end gap-2">
          {@render actions()}
        </div>
      {/if}
    </Dialog.Content>
  </Dialog.Portal>
</Dialog.Root>
