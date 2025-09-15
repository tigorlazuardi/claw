<script lang="ts">
  import type {
    AvailableSource,
    ListAvailableSourcesResponse,
  } from "../../gen/claw/v1/source_service_pb";
  import { useQuery } from "@sveltestack/svelte-query";
  import { getSourceServiceClient } from "../../connectrpc";

  interface Props {
    /**
     * Callback when data is loaded. If there is only one source, it will be auto-selected and passed as the second argument.l
     */
    onData?: (
      data: ListAvailableSourcesResponse,
      autoSelect?: AvailableSource,
    ) => void;
    /**
     * The currently selected source.
     */
    selected?: AvailableSource;
    /**
     * Whether the current selection is valid.
     */
    valid?: boolean;
    /**
     * The current value of the select input. Points to selected.name.
     */
    value?: string;
  }

  let {
    onData,
    selected = $bindable(),
    valid = $bindable(false),
    value = $bindable(""),
  }: Props = $props();

  const listAvailableSources = useQuery(
    ["sources", "add", "listDropDown"],
    () =>
      getSourceServiceClient().then((client) =>
        client.listAvailableSources({}),
      ),
    {
      onSuccess(data) {
        if (data.sources.length === 1) {
          selected = data.sources[0];
          value = data.sources[0].name;
          valid = true;
          onData?.(data, selected);
          return;
        }
        onData?.(data);
      },
    },
  );
  const allOk = $derived(valid && !!selected);
</script>

<fieldset class="fieldset">
  <legend class="fieldset-legend">
    <span
      class={{
        "text-success": allOk,
      }}
    >
      Source <span class="text-error">*</span>
    </span>
    {#if selected?.description}
      <div></div>
    {/if}
  </legend>
  {#if $listAvailableSources.isLoading}
    {@render loadingSources()}
  {:else if $listAvailableSources.isSuccess}
    {@render sourcesInput($listAvailableSources.data)}
  {:else}
    {@render sourcesError($listAvailableSources.error)}
  {/if}
</fieldset>

{#snippet loadingSources()}
  <select class="select w-full">
    <option disabled selected value="" class="loading loading-spinner"></option>
  </select>
  <span class="label">Getting list of sources. Please wait...</span>
{/snippet}

{#snippet sourcesInput(data: ListAvailableSourcesResponse)}
  {@const sources = data.sources}
  <select
    onchange={(e) => {
      valid = e.currentTarget.validity.valid;
      selected = $listAvailableSources.data?.sources.find(
        (s) => s.name === e.currentTarget.value,
      );
    }}
    bind:value
    class={{
      "select w-full": true,
      "select-success": allOk,
      "text-success": allOk,
    }}
    required
  >
    {#if !selected}
      <option disabled value="" class="text-base-100">
        -- select a source --
      </option>
    {/if}
    {#each sources as source (source.name)}
      <option class="text-base-content" value={source.name}>
        {source.displayName} ({source.name})
      </option>
    {/each}
  </select>
  <span
    class={{
      label: true,
      "text-success": allOk,
    }}
  >
    Choose supported source
  </span>
{/snippet}

{#snippet sourcesError(err: any)}
  <select class="select">
    <option disabled selected value="">-- error loading sources --</option>
  </select>
  <span class="label text-error">
    Error loading sources: {err}
  </span>
{/snippet}
