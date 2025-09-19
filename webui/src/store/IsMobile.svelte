<script lang="ts">
  import { onMount, type Snippet } from "svelte";
  import { isMobile } from "#/store/isMobile";

  interface Props {
    children: Snippet;
  }

  const { children }: Props = $props();

  function setIsMobile() {
    isMobile.set(window.innerWidth <= 640);
  }

  onMount(setIsMobile);
</script>

<svelte:window onresize={setIsMobile} />

{@render children()}
