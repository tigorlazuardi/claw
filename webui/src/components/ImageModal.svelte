<script lang="ts">
  import type { Image } from "../gen/claw/v1/image_pb";
  import type { M } from "../types";

  interface Props {
    image: M<Image>;
    onCloseRequest?: (event: Event) => void;
  }

  const { image, onCloseRequest }: Props = $props();

  function formatFileSize(bytes: number): string {
    const sizes = ["B", "KB", "MB", "GB"];
    if (bytes === 0) return "0 B";
    const i = Math.floor(Math.log(bytes) / Math.log(1024));
    return Math.round((bytes / Math.pow(1024, i)) * 10) / 10 + " " + sizes[i];
  }

  function formatDimensions(width: number, height: number): string {
    return `${width} × ${height}`;
  }
  function handleKeydown(event: KeyboardEvent) {
    if (event.key === "Escape") {
      onCloseRequest?.(event);
    }
  }
</script>

<div
  class="modal-overlay"
  onclick={onCloseRequest}
  onkeydown={handleKeydown}
  role="dialog"
  aria-modal="true"
  aria-labelledby="modal-title"
  tabindex="-1"
>
  <div class="modal-content" role="document">
    <button class="close-btn" onclick={onCloseRequest} aria-label="Close modal"
      >×</button
    >

    <div class="modal-image">
      <img
        src="https://picsum.photos/800/600?random={image.id}"
        alt={image.title || `Image ${image.id}`}
      />
    </div>

    <div class="modal-info">
      <h3 id="modal-title">{image.title || `Image ${image.id}`}</h3>
      <div class="info-grid">
        <div class="info-item">
          <span class="info-label">Dimensions:</span>
          <span class="info-value"
            >{formatDimensions(image.width, image.height)}</span
          >
        </div>
        <div class="info-item">
          <span class="info-label">File Size:</span>
          <span class="info-value">{formatFileSize(image.filesize)}</span>
        </div>
        {#if image.postAuthor}
          <div class="info-item">
            <span class="info-label">Author:</span>
            <span class="info-value">{image.postAuthor}</span>
          </div>
        {/if}
        {#if image.postUrl}
          <div class="info-item">
            <span class="info-label">Source:</span>
            <a
              href={image.postUrl}
              target="_blank"
              rel="noopener noreferrer"
              class="info-link"
            >
              View Original
            </a>
          </div>
        {/if}
      </div>

      {#if image.tags && image.tags.length > 0}
        <div class="tags-section">
          <span class="tags-label">Tags:</span>
          <div class="tags">
            {#each image.tags as tag}
              <span class="tag">{tag}</span>
            {/each}
          </div>
        </div>
      {/if}
    </div>
  </div>
</div>

<style>
  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: rgba(0, 0, 0, 0.8);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    backdrop-filter: blur(4px);
  }

  .modal-content {
    background-color: #2d2d2d;
    border-radius: 8px;
    padding: 1.5rem;
    max-width: 90vw;
    max-height: 90vh;
    overflow: auto;
    position: relative;
    display: flex;
    flex-direction: column;
    gap: 1rem;
    box-shadow: 0 10px 25px rgba(0, 0, 0, 0.5);
    pointer-events: all;
  }

  .close-btn {
    position: absolute;
    top: 1rem;
    right: 1rem;
    background: rgba(0, 0, 0, 0.6);
    border: none;
    color: #ffffff;
    font-size: 1.5rem;
    width: 32px;
    height: 32px;
    border-radius: 50%;
    cursor: pointer;
    z-index: 1001;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s ease;
  }

  .close-btn:hover {
    background-color: rgba(255, 255, 255, 0.2);
    transform: scale(1.1);
  }

  .modal-image {
    text-align: center;
    margin-bottom: 1rem;
  }

  .modal-image img {
    max-width: 100%;
    max-height: 70vh;
    object-fit: contain;
    border-radius: 4px;
  }

  .modal-info h3 {
    margin: 0 0 1rem 0;
    color: #ffffff;
    font-size: 1.25rem;
    font-weight: 600;
  }

  .info-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 0.75rem;
    margin-bottom: 1rem;
  }

  .info-item {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .info-label {
    font-size: 0.8rem;
    color: #aaa;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .info-value {
    color: #ffffff;
    font-size: 0.9rem;
  }

  .info-link {
    color: #646cff;
    text-decoration: none;
    font-size: 0.9rem;
    transition: color 0.2s ease;
  }

  .info-link:hover {
    color: #5a5fcf;
    text-decoration: underline;
  }

  .tags-section {
    border-top: 1px solid #3d3d3d;
    padding-top: 1rem;
  }

  .tags-label {
    font-size: 0.8rem;
    color: #aaa;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    display: block;
    margin-bottom: 0.5rem;
  }

  .tags {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
  }

  .tag {
    background-color: #3d3d3d;
    color: #ffffff;
    padding: 0.25rem 0.5rem;
    border-radius: 4px;
    font-size: 0.8rem;
    font-weight: 500;
  }

  @media (max-width: 768px) {
    .modal-content {
      margin: 1rem;
      padding: 1rem;
    }

    .info-grid {
      grid-template-columns: 1fr;
    }

    .modal-image img {
      max-height: 50vh;
    }
  }
</style>

