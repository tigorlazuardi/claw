<script lang="ts">
  import { QueryClient, QueryClientProvider } from "@sveltestack/svelte-query";
  import LoadingModal from "../components/LoadingModal.svelte";
  const queryClient = new QueryClient();

  let showAddModal = $state(false);
</script>

<QueryClientProvider client={queryClient}>
  <div class="p-[2rem] 2xl:max-w-[60vw] max-w-full m-auto">
    <div class="flex justify-between items-start mb-[2rem]">
      <div>
        <h1 class="text-4xl text-base-content">Sources</h1>
        <p class="text-lg font-light text-base-content/70">
          Configure and manage your image sources
        </p>
      </div>
      <button
        class="btn btn-lg btn-primary"
        onclick={() => (showAddModal = true)}
      >
        + Add
      </button>
    </div>

    <div class="sources-grid">
      <div class="source-card">
        <div class="source-header">
          <h3>🔴 Reddit</h3>
          <span class="source-status active">Active</span>
        </div>
        <div class="source-info">
          <p>
            <strong>Subreddits:</strong>
            wallpapers, EarthPorn, SkyPorn
          </p>
          <p>
            <strong>Schedule:</strong>
            Every 6 hours
          </p>
          <p>
            <strong>Last Run:</strong>
            2 hours ago
          </p>
          <p>
            <strong>Images Collected:</strong>
            1,247
          </p>
        </div>
        <div class="source-actions">
          <button class="btn-primary">Configure</button>
          <button class="btn-secondary">Run Now</button>
        </div>
      </div>

      <div class="source-card">
        <div class="source-header">
          <h3>🎨 Danbooru</h3>
          <span class="source-status active">Active</span>
        </div>
        <div class="source-info">
          <p>
            <strong>Tags:</strong>
            landscape, scenery, nature
          </p>
          <p>
            <strong>Schedule:</strong>
            Daily at 3:00 AM
          </p>
          <p>
            <strong>Last Run:</strong>
            8 hours ago
          </p>
          <p>
            <strong>Images Collected:</strong>
            432
          </p>
        </div>
        <div class="source-actions">
          <button class="btn-primary">Configure</button>
          <button class="btn-secondary">Run Now</button>
        </div>
      </div>

      <div class="source-card">
        <div class="source-header">
          <h3>📸 Unsplash</h3>
          <span class="source-status paused">Paused</span>
        </div>
        <div class="source-info">
          <p>
            <strong>Keywords:</strong>
            minimal, abstract, architecture
          </p>
          <p>
            <strong>Schedule:</strong>
            Every 12 hours
          </p>
          <p>
            <strong>Last Run:</strong>
            3 days ago
          </p>
          <p>
            <strong>Images Collected:</strong>
            89
          </p>
        </div>
        <div class="source-actions">
          <button class="btn-primary">Configure</button>
          <button class="btn-secondary">Resume</button>
        </div>
      </div>

      <div class="source-card">
        <div class="source-header">
          <h3>🌐 Custom RSS</h3>
          <span class="source-status inactive">Inactive</span>
        </div>
        <div class="source-info">
          <p>
            <strong>Feed URL:</strong>
            example.com/wallpapers.rss
          </p>
          <p>
            <strong>Schedule:</strong>
            Not configured
          </p>
          <p>
            <strong>Last Run:</strong>
            Never
          </p>
          <p>
            <strong>Images Collected:</strong>
            0
          </p>
        </div>
        <div class="source-actions">
          <button class="btn-primary">Configure</button>
          <button class="btn-secondary">Enable</button>
        </div>
      </div>
    </div>
  </div>

  {#if showAddModal}
    {#await import("../components/AddSourceModal.svelte")}
      <LoadingModal />
    {:then { default: AddSourceModal }}
      <AddSourceModal onCloseRequest={() => (showAddModal = false)} />
    {/await}
  {/if}
</QueryClientProvider>

<style>
  .sources-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
    gap: 1.5rem;
  }

  .source-card {
    background-color: hsl(0, 0%, 18%);
    border-radius: 8px;
    padding: 1.5rem;
    border-left: 4px solid hsl(235, 100%, 65%);
  }

  .source-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
  }

  .source-header h3 {
    margin: 0;
    color: hsl(0, 0%, 100%);
    font-size: 1.1rem;
    font-weight: 600;
  }

  .source-status {
    padding: 0.25rem 0.75rem;
    border-radius: 20px;
    font-size: 0.8rem;
    font-weight: 600;
    text-transform: uppercase;
  }

  .source-status.active {
    background-color: hsla(142, 76%, 36%, 0.2);
    color: hsl(142, 76%, 56%);
  }

  .source-status.paused {
    background-color: hsla(43, 96%, 56%, 0.2);
    color: hsl(43, 96%, 56%);
  }

  .source-status.inactive {
    background-color: hsla(0, 84%, 60%, 0.2);
    color: hsl(0, 84%, 60%);
  }

  .source-info {
    margin-bottom: 1.5rem;
  }

  .source-info p {
    margin: 0.5rem 0;
    color: hsl(0, 0%, 67%);
    font-size: 0.9rem;
  }

  .source-info strong {
    color: hsl(0, 0%, 100%);
  }

  .source-actions {
    display: flex;
    gap: 0.75rem;
  }

  .btn-primary,
  .btn-secondary {
    padding: 0.5rem 1rem;
    border-radius: 6px;
    border: none;
    font-size: 0.9rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .btn-primary {
    background-color: hsl(235, 100%, 65%);
    color: hsl(0, 0%, 100%);
  }

  .btn-primary:hover {
    background-color: hsl(235, 85%, 60%);
  }

  .btn-secondary {
    background-color: hsl(0, 0%, 24%);
    color: hsl(0, 0%, 67%);
    border: 1px solid hsl(0, 0%, 37%);
  }

  .btn-secondary:hover {
    background-color: hsl(0, 0%, 30%);
    color: hsl(0, 0%, 100%);
  }
</style>
