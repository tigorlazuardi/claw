<script lang="ts">
  import type { Image } from "../gen/claw/v1/image_pb";
  import type { M } from "../types";

  interface Props {
    image: M<Image>;
  }

  const { image }: Props = $props();

  function formatFileSize(bytes: number): string {
    const sizes = ["B", "KB", "MB", "GB"];
    if (bytes === 0) return "0 B";
    const i = Math.floor(Math.log(bytes) / Math.log(1024));
    return Math.round((bytes / Math.pow(1024, i)) * 10) / 10 + " " + sizes[i];
  }

  function formatDimensions(width: number, height: number): string {
    return `${width} × ${height}`;
  }

  let modalOpen = $state(false);
</script>

{#snippet urlify(title: string, link?: string)}
  {#if link}
    <a
      class="image-title"
      href={link}
      target="_blank"
      rel="noopener noreferrer"
      aria-label="(opens in a new tab)">{title}</a
    >
  {:else}
    <div class="image-title">{title}</div>
  {/if}
{/snippet}

<div>
  <figure class="image-card">
    <button
      class="image-container"
      tabindex="0"
      onclick={() => (modalOpen = true)}
    >
      <img
        src="https://picsum.photos/{image.width}/{image.height}?random={image.id}"
        alt={image.title || `Image ${image.id}`}
        loading="lazy"
      />
      <div class="image-overlay">
        <div class="image-info">
          {@render urlify(image.title || `Image ${image.id}`, image.postUrl)}
          <div class="image-meta">
            {formatDimensions(image.width, image.height)} • {formatFileSize(
              image.filesize,
            )}
          </div>
        </div>
      </div>
      {#if image.isFavorite}
        <div class="favorite-badge">★</div>
      {/if}
    </button>
  </figure>
  {#if modalOpen}
    {#await import("./ImageModal.svelte") then { default: ImageModal }}
      <ImageModal {image} onCloseRequest={() => (modalOpen = false)} />
    {/await}
  {/if}
</div>

<style>
  .image-card {
    border-radius: 0.25rem;
    overflow: hidden;
    cursor: pointer;
    transition:
      transform 0.2s ease,
      box-shadow 0.2s ease;
    outline: none;
  }

  .image-card:hover,
  .image-card:focus {
    transform: translateY(-2px);
  }

  .image-container {
    position: relative;
    overflow: hidden;
    min-height: 5rem; /* If image is broken, the title text should be at least readable*/
    margin: 0;
    padding: 0;
    text-align: left;
  }

  .image-container img {
    width: 100%;
    display: block;
    object-fit: contain;
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
    background-color: hsla(51, 100%, 50%, 0.9);
    color: hsl(0, 0%, 0%);
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
    background: linear-gradient(
      to top,
      hsla(0, 0%, 0%, 0.8) 0%,
      hsla(0, 0%, 0%, 0.4) 50%,
      transparent 100%
    );
    padding: 1rem 0.75rem 0.75rem 0.75rem;
    z-index: 2;
  }

  .image-title {
    font-weight: 600;
    color: hsl(0, 0%, 100%);
    margin-bottom: 0.25rem;
    font-size: 0.9rem;
    white-space: nowrap;
    overflow: hidden;
    text-align: left;
    text-overflow: ellipsis;
    text-shadow: 0 1px 2px hsla(0, 0%, 0%, 0.8);
  }

  .image-title:hover {
    text-decoration: underline;
    color: hsl(210, 100%, 80%);
  }

  .image-meta {
    color: hsl(0, 0%, 88%);
    font-size: 0.8rem;
    text-shadow: 0 1px 2px hsla(0, 0%, 0%, 0.8);
  }
</style>
