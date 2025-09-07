<script lang="ts">
  import Markdown from "svelte-exmarkdown";
  import { gfmPlugin } from "svelte-exmarkdown/gfm";
  const plugins = [gfmPlugin()];

  interface Props {
    text: string;
  }
  const { text }: Props = $props();

  const md = $state(text);
</script>

<Markdown {md} {plugins} />
