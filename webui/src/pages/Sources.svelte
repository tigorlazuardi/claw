<script lang="ts">
  import { useQuery } from "@sveltestack/svelte-query";
  import LoadingModal from "../components/LoadingModal.svelte";
  import { Button } from "flowbite-svelte";
  import { Hr } from "flowbite-svelte";
  import { getSourceServiceClient } from "../connectrpc";
  import type {
    ListSourcesRequest,
    ListSourcesResponse,
  } from "../gen/claw/v1/source_service_pb";
  import type { M } from "../types";
  import { toQuery, fromQuery } from "query-string-parser";
  import { Spinner } from "flowbite-svelte";

  let showAddModal = $state(false);
  type QueryState = Partial<M<ListSourcesRequest>>;

  let queryState: QueryState = $state(
    fromQuery(location.search.substring(1)) ?? {},
  );
  $effect(() => {
    const qs = toQuery(queryState);
    const search = qs ? `?${qs}` : "";
    if (location.search !== search) {
      history.replaceState(null, "", location.pathname + search);
    }
  });

  const listResult = useQuery(["sources", "list", queryState], async () => {
    const client = await getSourceServiceClient();
    return client.listSources(queryState);
  });
</script>

<div class="p-[2rem] 2xl:max-w-[70vw] max-w-full m-auto">
  <div class="flex justify-between items-start mb-[2rem]">
    <div>
      <h1 class="text-4xl">Sources</h1>
      <p class="text-lg font-light text-base-content/70">
        Configure and manage your image sources
      </p>
    </div>
    <Button onclick={() => (showAddModal = true)}>+ Add</Button>
  </div>
  <Hr />

  {#if $listResult.isLoading}
    {@render loading()}
  {:else if $listResult.isError}
    {@render error($listResult.error as Error)}
  {:else if $listResult.data}
    {@render data($listResult.data)}
  {/if}
</div>

{#if showAddModal}
  {#await import("../surfaces/AddSource.svelte")}
    <LoadingModal />
  {:then { default: AddSourceModal }}
    <AddSourceModal onCloseRequest={() => (showAddModal = false)} />
  {/await}
{/if}

{#snippet loading()}
  <div class="h-full w-full flex justify-center items-center">
    <Spinner />
  </div>
{/snippet}

{#snippet data(data: ListSourcesResponse)}
  <div></div>
{/snippet}

{#snippet error(error: Error)}
  <div class="p-[2rem] 2xl:max-w-[70vw] max-w-full m-auto">
    <div class="alert alert-error shadow-lg">
      <div>
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="stroke-current flex-shrink-0 h-6 w-6"
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
    </div>
  </div>
{/snippet}
