<script lang="ts">
  import LoadingModal from "../components/LoadingModal.svelte";
  import { getSourceServiceClient } from "../connectrpc";
  import type {
    ListSourcesRequest,
    ListSourcesResponse,
  } from "../gen/claw/v1/source_service_pb";
  import type { Source } from "../gen/claw/v1/source_pb";
  import IconImport from "@lucide/svelte/icons/import";
  import { resource } from "runed";

  import { toQuery, fromQuery } from "query-string-parser";
  import type { M } from "../types";

  let showAddModal = $state(false);

  type QueryState = Partial<M<ListSourcesRequest>>;

  let queryState: QueryState = $state(fromQuery(location.search) || {});

  $effect(() => {
    let qs = toQuery(queryState);
    if (qs !== "") {
      qs = `?${qs}`;
    }
    history.replaceState(null, "", location.pathname + qs);
  });

  const listResource = resource(
    () => queryState,
    async (queryState, _, { signal }) => {
      const client = await getSourceServiceClient({ signal });
      return client.listSources(queryState);
    },
    {
      debounce: 300,
    },
  );
</script>

<svelte:window onfocus={() => listResource.refetch()} />

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
  <div class="divider"></div>

  {#if listResource.loading}
    {@render loading()}
  {:else if listResource.error}
    {@render error(listResource.error)}
  {:else if listResource.current?.sources.length}
    {@render data(listResource.current)}
  {:else}
    {@render emptyState()}
  {/if}

  {#if showAddModal}
    {#await import("../surfaces/AddSource/AddSource.svelte")}
      <LoadingModal />
    {:then { default: AddSourceModal }}
      <AddSourceModal onCloseRequest={() => (showAddModal = false)} />
    {/await}
  {/if}
</div>

{#snippet loading()}
  <div class="flex justify-center items-center py-16">
    <span class="loading loading-spinner loading-lg"></span>
  </div>
{/snippet}

{#snippet data(data: ListSourcesResponse)}
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
    {#each data.sources as source}
      {@render sourceCard(source)}
    {/each}
  </div>
{/snippet}

{#snippet emptyState()}
  <div class="flex flex-col items-center justify-center py-16 text-center">
    <div class="mb-6">
      <IconImport class="mx-auto w-[6rem] h-[6rem] text-base-content" />
      <h3 class="text-xl font-semibold text-base-content mb-2">
        No Sources Have Been Added
      </h3>
      <p class="text-base-content/70 max-w-md">
        Get started by adding your first image source to begin collecting
        wallpapers and images.
      </p>
    </div>
    <button class="btn btn-primary" onclick={() => (showAddModal = true)}>
      <svg
        class="w-4 h-4"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M12 6v6m0 0v6m0-6h6m-6 0H6"
        ></path>
      </svg>
      Add Source
    </button>
  </div>
{/snippet}

{#snippet sourceCard(source: Source)}
  <div class="card bg-base-100 shadow-xl">
    <div class="card-body">
      <div class="flex justify-between items-start mb-2">
        <h2 class="card-title text-lg">
          {source.displayName}
        </h2>
        <div class="flex gap-2">
          {#if source.isDisabled}
            <div class="badge badge-error">Disabled</div>
          {:else}
            <div class="badge badge-success">Active</div>
          {/if}
          <div class="badge badge-info">{source.name}</div>
        </div>
      </div>

      <p class="text-sm text-base-content/70 mb-2">
        <strong>Parameter:</strong>
        {source.parameter || "None"}
      </p>

      <p class="text-sm text-base-content/70 mb-3">
        <strong>Countback:</strong>
        {source.countback}
      </p>

      {#if source.schedules.length > 0}
        <div class="mb-3">
          <p class="font-semibold text-sm mb-1">Schedules:</p>
          <div class="flex flex-wrap gap-1">
            {#each source.schedules as schedule}
              <div class="badge badge-secondary">{schedule.schedule}</div>
            {/each}
          </div>
        </div>
      {/if}

      <div class="text-xs text-base-content/60">
        {#if source.lastRunAt}
          <p>
            Last run: {new Date(
              Number(source.lastRunAt.seconds) * 1000,
            ).toLocaleString()}
          </p>
        {/if}
        <p>
          Created: {source.createdAt
            ? new Date(Number(source.createdAt.seconds) * 1000).toLocaleString()
            : "Unknown"}
        </p>
      </div>
    </div>
  </div>
{/snippet}

{#snippet error(error: Error)}
  <div class="alert alert-error">
    <svg
      xmlns="http://www.w3.org/2000/svg"
      class="stroke-current shrink-0 h-6 w-6"
      fill="none"
      viewBox="0 0 24 24"
    >
      <path
        stroke-linecap="round"
        stroke-linejoin="round"
        stroke-width="2"
        d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
      />
    </svg>
    <span>{error.message}</span>
  </div>
{/snippet}
