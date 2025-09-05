<script lang="ts">
  import type { Image } from '../gen/claw/v1/image_pb'
  import ImageCard from '../components/ImageCard.svelte'
  import { createEventDispatcher } from 'svelte'

  export let images: Image[] = []
  export let dateHeader: string = ''
  
  const dispatch = createEventDispatcher()
  
  function handleImageClick(event) {
    dispatch('imageClick', event.detail)
  }
</script>

<div class="date-group">
  {#if dateHeader}
    <h2 class="date-header">{dateHeader}</h2>
  {/if}
  
  <div class="masonry-grid">
    {#each images as image (image.id)}
      <ImageCard 
        {image} 
        on:click={handleImageClick}
      />
    {/each}
  </div>
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

  .masonry-grid {
    display: grid;
    grid-template-columns: repeat(6, 1fr);
    grid-auto-rows: 100px;
    gap: 1rem;
  }
</style>