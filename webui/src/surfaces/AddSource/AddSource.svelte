<script lang="ts">
  import {
    useMutation,
    useQuery,
    useQueryClient,
  } from "@sveltestack/svelte-query";
  import type { M } from "../../types";
  import type {
    AvailableSource,
    CreateSourceRequest,
    ValidateSourceParametersRequest,
  } from "../../gen/claw/v1/source_service_pb";
  import IconX from "@lucide/svelte/icons/x";
  import IconInfo from "@lucide/svelte/icons/info";
  import { getSourceServiceClient } from "../../connectrpc";
  import { Popover } from "bits-ui";
  import SelectSource from "./SelectSource.svelte";

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

  type Schedule = {
    pattern: string;
    nextRun?: Date;
  };
  const schedules = $state<Schedule[]>([]);
  const schedulePatterns = $derived(schedules.map((s) => s.pattern));

  let addSourceRequest = $state<Omit<M<CreateSourceRequest>, "schedules">>({
    name: "",
    parameter: "",
    displayName: "",
    countback: 0,
    isDisabled: false,
  });

  let selectedSource: AvailableSource | undefined = $state();
  let nameValid: boolean = $state(false);
  let parameterValid: boolean = $state(false);

  let canSubmitForm = $derived.by(() => {
    return [
      addSourceRequest.name,
      nameValid,
      selectedSource?.requireParameter && addSourceRequest.parameter,
      addSourceRequest.displayName,
      parameterValid,
    ].every((f) => f);
  });

  let parameterInputField: HTMLTextAreaElement | undefined = $state();
  let validateParameterRequest = $derived<M<ValidateSourceParametersRequest>>({
    sourceName: addSourceRequest.name,
    parameter: addSourceRequest.parameter,
  });

  const validateParameterQuery = useQuery(
    ["source", "add", "validateParameter"],
    () =>
      getSourceServiceClient().then((client) =>
        client.validateSourceParameters(validateParameterRequest),
      ),
    {
      onSuccess(data) {
        parameterInputField?.setCustomValidity("");
        if (data.transformedParameter) {
          addSourceRequest.parameter = data.transformedParameter;
        }
      },
      onError(err: Error) {
        parameterInputField?.setCustomValidity(err.message);
      },
      enabled: false, // Only run on demand
    },
  );

  function handleValidateParamaterOnBlur() {
    if (addSourceRequest.name && addSourceRequest.parameter) {
      useQueryClient().cancelQueries(["source", "add", "validateParameter"]);
      $validateParameterQuery.refetch();
    }
  }

  const { onCloseRequest }: Props = $props();

  const createSourceMutation = useMutation((req: M<CreateSourceRequest>) =>
    getSourceServiceClient().then((client) => client.createSource(req)),
  );

  function handleOnSubmit() {
    if (!canSubmitForm) {
      return;
    }
    $createSourceMutation.mutate(
      {
        ...addSourceRequest,
        schedules: schedulePatterns,
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

  let hasParameterHelp = $derived(
    selectedSource?.parameterHelp && selectedSource.parameterHelp.trim() !== "",
  );

  let scheduleInput: HTMLInputElement | undefined = $state();

  async function addSchedule() {}
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
      <SelectSource
        bind:selected={selectedSource}
        bind:valid={nameValid}
        bind:value={addSourceRequest.name}
      />
      {#if selectedSource}
        {#await import("./ParameterInput.svelte") then { default: ParameterInput }}
          <ParameterInput
            source={selectedSource}
            bind:value={addSourceRequest.parameter}
            bind:valid={parameterValid}
          />
        {/await}
        {@render displayNameInput()}
        {@render countbackInput(selectedSource)}
        {@render scheduleInputField()}
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
      bind:value={addSourceRequest.displayName}
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
      bind:value={addSourceRequest.countback}
    />
    <p class="label">
      The number of items to look up for. If value is 0, source default will be
      used
    </p>
  </fieldset>
{/snippet}

{#snippet scheduleInputField()}
  <fieldset class="fieldset">
    <legend class="fieldset-legend">
      <span>Schedules</span>
    </legend>
    <input
      type="text"
      class="input w-full"
      bind:value={addSourceRequest.displayName}
      placeholder="Schedule pattern"
      required
    />
    <p class="label"></p>
  </fieldset>
{/snippet}
