<script lang="ts">
  import { useQuery, useQueryClient } from "@sveltestack/svelte-query";
  import { form as formValidator, field } from "svelte-forms";
  import { required } from "svelte-forms/validators";
  import type { M } from "../types";
  import type {
    AvailableSource,
    CreateSourceRequest,
    ListAvailableSourcesResponse,
  } from "../gen/claw/v1/source_service_pb";
  import IconX from "@lucide/svelte/icons/x";
  import IconInfo from "@lucide/svelte/icons/info";
  import { getSourceServiceClient } from "../connectrpc";

  interface Props {
    /**
     * Callback when the modal is requesting to be closed or removed from the DOM.
     * This will be called when the user clicks the close button, cancel button,
     * or the server responds with success after form submission.
     *
     * This function will not be called if the user clicks outside the modal or press the escape key
     * to prevent accidental closure.
     */
    onCloseRequest: () => void;
  }

  const name = field("name", "", [required()]);
  const parameter = field("parameter", "");
  const displayName = field("display_name", "", [required()]);
  const countback = field("countback", 0);
  const isDisabled = field("is_disabled", false);
  const runOnCreation = field("run_on_creation", false);

  let scheduleExpressionError = $state<string | undefined>();
  const scheduleInput = field("schedule", "", [
    async (val: string) => {
      if (val.trim() === "") {
        scheduleExpressionError = undefined;
        return { valid: true, name: "cron_format" };
      }
      const { CronExpressionParser } = await import("cron-parser");
      try {
        CronExpressionParser.parse(val);
        scheduleExpressionError = undefined;
        return { valid: true, name: "cron_format" };
      } catch (e) {
        scheduleExpressionError =
          "invalid cron expression: " + (e as Error).message;
        return { valid: false, name: "cron_format" };
      }
    },
  ]);
  type Schedule = {
    pattern: string;
    nextRun?: Date;
  };
  const schedules = $state<Schedule[]>([]);
  const addSourceForm = formValidator(
    name,
    parameter,
    displayName,
    countback,
    isDisabled,
    scheduleInput,
    runOnCreation,
  );

  const { onCloseRequest }: Props = $props();

  async function postCreateSource() {
    if (!$addSourceForm.valid) {
      return;
    }
    const data: M<CreateSourceRequest> = {
      name: $name.value,
      parameter: $parameter.value,
      displayName: $displayName.value,
      countback: $countback.value,
      isDisabled: $isDisabled.value,
      schedules: schedules.map((s) => s.pattern),
    };
    return getSourceServiceClient().then((client) => client.createSource(data));
  }

  let serverErrorResponse = $state<Error | null>(null);
  let isSubmitting = $state(false);
  async function handleOnSubmit() {
    isSubmitting = true;
    return postCreateSource()
      .then(async (res) => {
        if (!res) return;
        const queryClient = useQueryClient();
        queryClient.invalidateQueries(["sources", "list"]);
        onCloseRequest();
      })
      .catch((err: Error) => (serverErrorResponse = err))
      .finally(() => (isSubmitting = false));
  }
  let showParameterHelp = $state(false);
  async function listSources() {
    return getSourceServiceClient().then((client) =>
      client.listAvailableSources({}),
    );
  }

  const listSourcesResult = useQuery(["sources", "listDropdown"], listSources);
  let selectedSource = $derived.by(() => {
    const sources = $listSourcesResult.data?.sources;
    if (sources?.length === 1) {
      $name.value = sources[0].name;
      return sources[0];
    }
    return sources?.find((s) => s.name === $name.value);
  });
  let hasParameterHelp = $derived(
    selectedSource?.parameterHelp && selectedSource.parameterHelp.trim() !== "",
  );

  async function addSchedule() {
    await scheduleInput.validate();
    if ($addSourceForm.hasError("schedule.cron_format")) return;
    const { CronExpressionParser } = await import("cron-parser");
    schedules.push({
      pattern: $scheduleInput.value,
      nextRun: $scheduleInput.value
        ? CronExpressionParser.parse($scheduleInput.value).next().toDate()
        : undefined,
    });
    scheduleInput.reset();
  }
</script>

<div class="modal modal-open">
  <div class="modal-box xl:w-[60vw] xl:max-w-[60vw]">
    {@render modalHeader()}
    <div class="divider"></div>
    <form
      class="w-full"
      onsubmit={(e) => {
        e.preventDefault();
        handleOnSubmit();
      }}
    >
      {@render selectSourceInput()}
      {#if selectedSource}
        {@render parameterInput(selectedSource)}
      {/if}
    </form>
  </div>
</div>

{#snippet modalHeader()}
  <div id="modal-header" class="flex justify-between items-center mb-4">
    <h2 class="text-2xl">Add New Source</h2>
    <button
      class="btn btn-square btn-ghost"
      type="button"
      onclick={onCloseRequest}
    >
      <IconX />
    </button>
  </div>
{/snippet}

{#snippet selectSourceInput()}
  <fieldset class="fieldset">
    <legend class="fieldset-legend">
      <span>
        Source <span class="text-error">*</span>
      </span>
      {#if selectedSource?.description}
        <div></div>
      {/if}
    </legend>
    {#if $listSourcesResult.isLoading}
      {@render loadingSources()}
    {:else if $listSourcesResult.isSuccess}
      {@render sourcesInput($listSourcesResult.data)}
    {:else}
      {@render sourcesError($listSourcesResult.error)}
    {/if}
  </fieldset>
{/snippet}

{#snippet loadingSources()}
  <select class="select">
    <option disabled selected value="" class="loading loading-spinner"></option>
  </select>
  <span class="label">Getting list of sources. Please wait...</span>
{/snippet}

{#snippet sourcesInput(data: ListAvailableSourcesResponse)}
  {@const sources = data.sources}
  <select class="select w-full" bind:value={$name.value} required>
    {#if !selectedSource}
      <option disabled value="" class="text-base-100">
        -- select a source --
      </option>
    {/if}
    {#each sources as source (source.name)}
      <option value={source.name}>
        {source.displayName} ({source.name})
      </option>
    {/each}
  </select>
  <span class="label">Choose supported source</span>
{/snippet}

{#snippet sourcesError(err: any)}
  <select class="select">
    <option disabled selected value="">-- error loading sources --</option>
  </select>
  <span class="label text-error">
    Error loading sources: {err}
  </span>
{/snippet}

{#snippet parameterInput(source: AvailableSource)}
  <fieldset class="fieldset">
    <legend class="fieldset-legend">
      <span>Parameters</span>
      {#if source.requireParameter}
        <span class="text-error">*</span>
      {/if}
      {#if source.parameterHelp}
        <button
          class="btn btn-square btn-ghost btn-xs"
          type="button"
          onclick={() => (showParameterHelp = !showParameterHelp)}
        >
          <IconInfo />
        </button>
      {/if}
    </legend>
    <textarea
      class="textarea h-[3rem] w-full"
      placeholder={source.parameterPlaceholder ||
        "Configuration parameters (JSON, comma-separated values, etc.)"}
      bind:value={$parameter.value}
      required={source.requireParameter}
    ></textarea>
    {#if source.requireParameter}
      <span class="label">Required</span>
    {:else}
      <span class="label">Optional</span>
    {/if}
  </fieldset>
{/snippet}
