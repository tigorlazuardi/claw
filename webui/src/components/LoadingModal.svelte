<script lang="ts">
  import SpinningArrow from "./SpinningArrow.svelte";

  interface Props {
    message?: string;
    onCancel?: () => void;
  }

  const { message = "Loading...", onCancel }: Props = $props();
</script>

<div class="modal-overlay" role="button" tabindex="0" aria-label="Loading modal">
  <div
    class="modal-content"
    role="dialog"
    aria-modal="true"
    aria-label="Loading"
    tabindex="-1"
    onclick={(e) => e.stopPropagation()}
    onkeydown={(e) => e.stopPropagation()}
  >
    <div class="loading-content">
      <SpinningArrow size="3rem" />
      <p class="loading-message">{message}</p>
      {#if onCancel}
        <button type="button" class="cancel-btn" onclick={onCancel}>
          Cancel
        </button>
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
    background-color: hsla(0, 0%, 0%, 0.7);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    padding: 1rem;
  }

  .modal-content {
    background-color: hsl(0, 0%, 18%);
    border-radius: 12px;
    padding: 2rem;
    box-shadow: 0 25px 50px hsla(0, 0%, 0%, 0.5);
    min-width: 300px;
    max-width: 400px;
  }

  .loading-content {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1.5rem;
    text-align: center;
  }

  .loading-message {
    color: hsl(0, 0%, 100%);
    font-size: 1.1rem;
    margin: 0;
    font-weight: 500;
  }

  .cancel-btn {
    padding: 0.75rem 1.5rem;
    border: 1px solid hsl(0, 0%, 37%);
    border-radius: 6px;
    background-color: hsl(0, 0%, 24%);
    color: hsl(0, 0%, 67%);
    font-size: 0.9rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .cancel-btn:hover {
    background-color: hsl(0, 0%, 30%);
    color: hsl(0, 0%, 100%);
    border-color: hsl(0, 0%, 45%);
  }
</style>