<script lang="ts">
  import { createMutation, useQueryClient } from "@tanstack/svelte-query";
  import type { M } from "#/types";
  import type {
    AvailableSource,
    CreateSourceRequest,
  } from "../../../gen/claw/v1/source_service_pb";
  import { getSourceServiceClient } from "../../../connectrpc";
  import SelectSource from "./SelectSource.svelte";
  import DialogModal from "#/components/DialogModal.svelte";

  interface Props {
    open: boolean;
  }

  let addSourceRequest = $state<M<CreateSourceRequest>>({
    name: "",
    parameter: "",
    displayName: "",
    countback: 0,
    isDisabled: false,
    schedules: [],
  });

  const queryClient = useQueryClient();

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

  let { open = $bindable(false) }: Props = $props();

  const createSourceMutation = createMutation({
    mutationKey: ["sources", "create"],
    mutationFn: async function (req: M<CreateSourceRequest>) {
      const client = await getSourceServiceClient();
      return client.createSource(req);
    },
  });

  async function handleOnSubmit() {
    if (!canSubmitForm) {
      return;
    }
    return $createSourceMutation.mutateAsync(addSourceRequest, {
      onSuccess: () => {
        queryClient.invalidateQueries({
          queryKey: ["sources", "list"],
        });
        open = false;
      },
    });
  }
</script>

<DialogModal bind:open class="w-[90vw] sm:w-[60vw]">
  {#snippet title()}
    <span class="font-bold">Add New Source</span>
  {/snippet}
  {#snippet description()}
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
        {#await import("./SchedulesInput.svelte") then { default: SchedulesInput }}
          <SchedulesInput bind:value={addSourceRequest.schedules} />
        {/await}
        {#await import("./Actions.svelte") then { default: Actions }}
          <Actions
            onclick={(evt, immediate) => {
              // TODO: Handle immediate run after creation
              evt.preventDefault();
              handleOnSubmit();
            }}
            oncancel={(evt) => {
              evt.preventDefault();
              onCloseRequest();
            }}
            disabled={!canSubmitForm || $createSourceMutation.isPending}
          />
        {/await}
      {/if}
    </form>
  {/snippet}
</DialogModal>
