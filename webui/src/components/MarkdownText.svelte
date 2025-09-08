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

<div class="prose prose-li:leading-none">
  <Markdown {md} {plugins} />
</div>
