<script lang="ts">
  import type { Image } from '../gen/claw/v1/image_pb'
  import { createEventDispatcher } from 'svelte'

  export let image: Image
  
  const dispatch = createEventDispatcher()
  
  function formatFileSize(bytes: number): string {
    const sizes = ['B', 'KB', 'MB', 'GB']
    if (bytes === 0) return '0 B'
    const i = Math.floor(Math.log(bytes) / Math.log(1024))
    return Math.round(bytes / Math.pow(1024, i) * 10) / 10 + ' ' + sizes[i]
  }
  
  function formatDimensions(width: number, height: number): string {
    return `${width} × ${height}`
  }

  function handleClick() {
    dispatch('click', { image })
  }
</script>

<div 
  class="image-card" 
  style="grid-column: span {Math.ceil(image.width / image.height * 2)}; grid-row: span {Math.ceil(image.height / image.width * 2)}"
  on:click={handleClick}
  on:keydown={(e) => e.key === 'Enter' && handleClick()}
  role="button"
  tabindex="0"
>
  <div class="image-container">
    <img 
      src="https://picsum.photos/{image.width}/{image.height}?random={image.id}" 
      alt={image.title || `Image ${image.id}`}
      loading="lazy"
    />
    <div class="image-overlay">
      <div class="image-info">
        <div class="image-title">
          {image.title || `Image ${image.id}`}
        </div>
        <div class="image-meta">
          {formatDimensions(image.width, image.height)} • {formatFileSize(image.filesize)}
        </div>
      </div>
    </div>
    {#if image.isFavorite}
      <div class="favorite-badge">★</div>
    {/if}
  </div>
</div>

<style>
  .image-card {
    background-color: #2d2d2d;
    border-radius: 8px;
    overflow: hidden;
    cursor: pointer;
    transition: transform 0.2s ease, box-shadow 0.2s ease;
    outline: none;
  }
  
  .image-card:hover,
  .image-card:focus {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  }
  
  .image-container {
    position: relative;
    width: 100%;
    height: 100%;
    overflow: hidden;
  }
  
  .image-container img {
    width: 100%;
    height: 100%;
    display: block;
    object-fit: cover;
    transition: transform 0.2s ease;
  }
  
  .image-card:hover .image-container img,
  .image-card:focus .image-container img {
    transform: scale(1.05);
  }
  
  .favorite-badge {
    position: absolute;
    top: 0.5rem;
    right: 0.5rem;
    background-color: rgba(255, 215, 0, 0.9);
    color: #000;
    padding: 0.25rem;
    border-radius: 50%;
    font-size: 0.8rem;
    width: 1.5rem;
    height: 1.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 3;
  }
  
  .image-overlay {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    background: linear-gradient(to top, rgba(0, 0, 0, 0.8) 0%, rgba(0, 0, 0, 0.4) 50%, transparent 100%);
    padding: 1rem 0.75rem 0.75rem 0.75rem;
    z-index: 2;
  }
  
  .image-title {
    font-weight: 600;
    color: #ffffff;
    margin-bottom: 0.25rem;
    font-size: 0.9rem;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.8);
  }
  
  .image-meta {
    color: #e0e0e0;
    font-size: 0.8rem;
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.8);
  }
</style>