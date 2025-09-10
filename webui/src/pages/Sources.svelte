<script lang="ts">
  import { useQuery } from "@sveltestack/svelte-query";
  import LoadingModal from "../components/LoadingModal.svelte";
  import { Button } from "flowbite-svelte";
  import { Hr } from "flowbite-svelte";
  import { Card, Badge, P, Heading } from "flowbite-svelte";
  import { getSourceServiceClient } from "../connectrpc";
  import type {
    ListSourcesRequest,
    ListSourcesResponse,
  } from "../gen/claw/v1/source_service_pb";
  import type { M } from "../types";
  import { toQuery, fromQuery } from "query-string-parser";
  import { Spinner } from "flowbite-svelte";
  import IconImport from "@lucide/svelte/icons/import";

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
  {#if data.sources.length === 0}
    <div class="flex flex-col items-center justify-center py-16 text-center">
      <div class="mb-6">
        <IconImport
          class="mx-auto text-gray-500 dark:text-gray-300 w-[6rem] h-[6rem]"
        />
        <Heading
          tag="h3"
          class="text-xl font-semibold text-gray-700 dark:text-gray-300 mb-2"
        >
          No Sources Have Been Added
        </Heading>
        <P class="text-gray-500 dark:text-gray-400 max-w-md text-center">
          To start, create a source to begin collecting wallpapers and images
        </P>
      </div>
      <Button
        onclick={() => (showAddModal = true)}
        class="flex items-center gap-2"
      >
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
      </Button>
    </div>
  {:else}
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      {#each data.sources as source}
        {@render sourceCard(source)}
      {/each}
    </div>
  {/if}
{/snippet}

{#snippet sourceCard(source: import("../gen/claw/v1/source_pb").Source)}
  <Card class="max-w-none">
    <div class="flex justify-between items-start mb-2">
      <Heading tag="h5" class="mb-2 text-xl font-semibold tracking-tight">
        {source.displayName}
      </Heading>
      <div class="flex gap-2">
        {#if source.isDisabled}
          <Badge color="red">Disabled</Badge>
        {:else}
          <Badge color="green">Active</Badge>
        {/if}
        <Badge color="blue">{source.name}</Badge>
      </div>
    </div>

    <P class="mb-3 font-normal text-gray-700 dark:text-gray-400 leading-tight">
      <strong>Parameter:</strong>
      {source.parameter || "None"}
    </P>

    <P class="mb-3 font-normal text-gray-700 dark:text-gray-400 leading-tight">
      <strong>Countback:</strong>
      {source.countback}
    </P>

    {#if source.schedules.length > 0}
      <div class="mb-3">
        <P class="font-semibold text-sm mb-1">Schedules:</P>
        {#each source.schedules as schedule}
          <Badge color="purple" class="mr-1 mb-1">{schedule.schedule}</Badge>
        {/each}
      </div>
    {/if}

    <div class="text-xs text-gray-500">
      {#if source.lastRunAt}
        <P>
          Last run: {new Date(
            Number(source.lastRunAt.seconds) * 1000,
          ).toLocaleString()}
        </P>
      {/if}
      <P>
        Created: {source.createdAt
          ? new Date(Number(source.createdAt.seconds) * 1000).toLocaleString()
          : "Unknown"}
      </P>
    </div>
  </Card>
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
