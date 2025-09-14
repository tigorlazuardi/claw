<script lang="ts">
  import {
    useMutation,
    useQuery,
    useQueryClient,
  } from "@sveltestack/svelte-query";
  import { form as formValidator, field } from "svelte-forms";
  import { required } from "svelte-forms/validators";
  import type { M } from "../types";
  import type {
    AvailableSource,
    CreateSourceRequest,
    ListAvailableSourcesResponse,
    ValidateSourceParametersRequest,
  } from "../gen/claw/v1/source_service_pb";
  import IconX from "@lucide/svelte/icons/x";
  import IconInfo from "@lucide/svelte/icons/info";
  import { getSourceServiceClient } from "../connectrpc";
  import { Popover } from "bits-ui";

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

  const validateParameterMutation = useMutation(
    (req: M<ValidateSourceParametersRequest>) =>
      getSourceServiceClient().then((client) =>
        client.validateSourceParameters(req),
      ),
    {
      onSuccess(data) {
        if (data.transformedParameter) {
          $parameter.value = data.transformedParameter;
        }
      },
    },
  );

  const name = field("name", "", [required()]);
  const parameter = field("parameter", "", [
    async function (val: string) {
      if (val.trim() === "") {
        return {
          valid: !!selectedSource?.requireParameter,
          name: "parameter_value",
        };
      }
      try {
        await $validateParameterMutation.mutateAsync({
          sourceName: $name.value,
          parameter: val,
        });
        return { valid: true, name: "parameter_value" };
      } catch {
        return { valid: false, name: "parameter_value" };
      }
    },
  ]);
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

  const createSourceMutation = useMutation((req: M<CreateSourceRequest>) =>
    getSourceServiceClient().then((client) => client.createSource(req)),
  );

  function validatePareameter() {
    if (!selectedSource) return;
    if (!$parameter.valid) return;
    if (!$parameter.value) return;
    $validateParameterMutation.mutate({
      sourceName: $name.value,
      parameter: $parameter.value,
    });
  }

  function handleOnSubmit() {
    if (!$addSourceForm.valid) {
      return;
    }
    $createSourceMutation.mutate(
      {
        name: $name.value,
        parameter: $parameter.value,
        displayName: $displayName.value,
        countback: $countback.value,
        isDisabled: $isDisabled.value,
        schedules: schedules.map((s) => s.pattern),
      },
      {
        onSuccess: () => {
          useQueryClient().invalidateQueries(["sources", "list"]);
          onCloseRequest();
        },
      },
    );
  }
  let showParameterHelp = $state(false);

  const listAvailableSources = useQuery(
    ["sources", "add", "listDropDown"],
    () =>
      getSourceServiceClient().then((client) =>
        client.listAvailableSources({}),
      ),
    {
      onSuccess(data) {
        if (data.sources.length === 1) {
          if ($countback.value === 0) {
            $countback.value = data.sources[0].defaultCountback;
          }
          $name.value = data.sources[0].name;
        }
      },
    },
  );
  let selectedSource = $derived(
    $listAvailableSources.data?.sources.find((s) => s.name === $name.value),
  );

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
        {@render displayNameInput()}
        {@render countbackInput(selectedSource)}
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
    {#if $listAvailableSources.isLoading}
      {@render loadingSources()}
    {:else if $listAvailableSources.isSuccess}
      {@render sourcesInput($listAvailableSources.data)}
    {:else}
      {@render sourcesError($listAvailableSources.error)}
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
      {#if hasParameterHelp}
        <Popover.Root onOpenChange={(v) => (showParameterHelp = v)}>
          <Popover.Trigger
            class="btn btn-square btn-ghost btn-xs"
            type="button"
          >
            <IconInfo />
          </Popover.Trigger>
          <Popover.Portal>
            <Popover.Content
              class="z-[9999] bg-transparent"
              data-theme="dracula"
            >
              <div class="card bg-base-300 border-base-200 border">
                <div class="card-body">
                  <div class="card-title">Parameter Help</div>
                  <div
                    class="p-4 border border-base-100"
                    style="box-shadow: inset 0 6px 12px rgba(0, 0, 0, 0.15), inset 0 2px 4px rgba(0, 0, 0, 0.1);"
                  >
                    {#if showParameterHelp}
                      {#await import ("../components/MarkdownText.svelte") then { default: MarkdownText }}
                        <MarkdownText text={source.parameterHelp} />
                      {/await}
                    {/if}
                  </div>
                </div>
              </div>
            </Popover.Content>
          </Popover.Portal>
        </Popover.Root>
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

{#snippet displayNameInput()}
  <fieldset class="fieldset">
    <legend class="fieldset-legend">
      <span>
        Display Name <span class="text-error">*</span>
      </span>
    </legend>
    <input
      type="text"
      class="input w-full"
      bind:value={$displayName.value}
      placeholder="Human readable name for the UI"
      required
    />
    <p class="label">Required</p>
  </fieldset>
{/snippet}

{#snippet countbackInput(source: AvailableSource)}
  <fieldset class="fieldset">
    <legend class="fieldset-legend">
      <span>Count Back</span>
      <Popover.Root>
        <Popover.Trigger class="btn btn-square btn-ghost btn-xs" type="button">
          <IconInfo />
        </Popover.Trigger>
        <Popover.Portal>
          <Popover.Content class="z-[9999] bg-transparent" data-theme="dracula">
            <div class="card bg-base-300 border-base-200 border">
              <div class="card-body">
                <div class="card-title flex-col items-start">
                  <span>Count Back</span>
                  <span class="text-xs">
                    Souce: {source.displayName} ({source.name})
                  </span>
                  <span class="text-xs">
                    Default Value: {source.defaultCountback}
                  </span>
                </div>
                <div
                  class="prose p-4 border border-base-100"
                  style="box-shadow: inset 0 6px 12px rgba(0, 0, 0, 0.15), inset 0 2px 4px rgba(0, 0, 0, 0.1);"
                >
                  <p>
                    Count Back lookup the number of "Items" to look up for and
                    sets an upper limit for Claw to stop looking for more
                    images.
                  </p>
                  <p>
                    However, "Items" do not mean "Images". It can be something
                    like a post or forum entry, but if it's not an image, the
                    post will be skipped, however it still counts towards the
                    count back limit.
                  </p>
                  <p>
                    The main purpose for Count Back is to avoid API rate limit.
                  </p>

                  <p>
                    Source Default Value usually expects for the same Source
                    Kind to not have overlapping schedules and takes as many
                    images possible under the API limit.
                  </p>
                </div>
              </div>
            </div>
          </Popover.Content>
        </Popover.Portal>
      </Popover.Root>
    </legend>
    <input
      class="input w-full"
      type="number"
      step="1"
      min="0"
      bind:value={$countback.value}
    />
    <p class="label">
      The number of items to look up for. If value is 0, source default will be
      used
    </p>
  </fieldset>
{/snippet}
