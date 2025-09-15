<script lang="ts">
  import { useMutation, useQueryClient } from "@sveltestack/svelte-query";
  import type { M } from "../../types";
  import type {
    AvailableSource,
    CreateSourceRequest,
  } from "../../gen/claw/v1/source_service_pb";
  import IconX from "@lucide/svelte/icons/x";
  import { getSourceServiceClient } from "../../connectrpc";
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
        {#await import("./DisplayNameInput.svelte") then { default: DisplayNameInput }}
          <DisplayNameInput bind:value={addSourceRequest.displayName} />
        {/await}
        {#await import("./CountbackInput.svelte") then { default: CountbackInput }}
          <CountbackInput
            bind:value={addSourceRequest.countback}
            source={selectedSource}
          />
        {/await}
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
