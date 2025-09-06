<script lang="ts">
  import type { Image } from "../gen/claw/v1/image_pb";
  import ImageCard from "../components/ImageCard.svelte";
  import Masonry from "svelte-bricks";

  export let images: Image[] = [];
  export let dateHeader: string = "";
</script>

<div class="date-group">
  {#if dateHeader}
    <h2 class="date-header">{dateHeader}</h2>
  {/if}

  <Masonry items={images}>
    {#snippet children({ item })}
      <ImageCard image={item} />
    {/snippet}
  </Masonry>
</div>

<style>
  .date-group {
    margin-bottom: 3rem;
  }

  .date-header {
    font-size: 1.1rem;
    font-weight: 600;
    color: #ffffff;
    margin: 0 0 1rem 0;
    padding: 0.5rem;
    background-color: #2a2a2a;
    border-radius: 6px;
    border-left: 4px solid #646cff;
  }
</style>

