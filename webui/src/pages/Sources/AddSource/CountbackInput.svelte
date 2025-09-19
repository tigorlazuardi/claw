<script lang="ts">
  import type { AvailableSource } from "#/gen/claw/v1/source_service_pb";
  interface Props {
    value: number;
    source: AvailableSource;
  }
  import PopoverInfo from "#/components/PopoverInfo.svelte";

  let { value = $bindable(0), source }: Props = $props();

  const allOk = $derived(value > 0);
</script>

<fieldset class="fieldset">
  <legend class="fieldset-legend">
    <span class={{ "text-success": allOk }}>Count Back</span>
    <PopoverInfo>
      {#snippet title()}
        <div class="card-title flex-col items-start">
          <span>Count Back</span>
          <span class="text-xs">
            Souce: {source.displayName} ({source.name})
          </span>
          <span class="text-xs">
            Default Value: {source.defaultCountback}
          </span>
        </div>
      {/snippet}

      <div class="prose text-wrap">
        <p>
          Count Back lookup the number of "Items" to look up for and sets an
          upper limit for Claw to stop looking for more images.
        </p>
        <p>
          However, "Items" do not mean "Images". It can be something like a post
          or forum entry, but if it's not an image, the post will be skipped,
          however it still counts towards the count back limit.
        </p>
        <p>
          The main purpose for Count Back is to avoid API rate limit and avoid
          Claw from doing too many requests to a Source and cause degradation in
          performance
        </p>

        <p>
          Source Default Value usually expects for the same Source Kind to not
          have overlapping schedules and takes as many images possible under the
          API limit.
        </p>
      </div>
    </PopoverInfo>
  </legend>
  <input
    class={{
      "input w-full": true,
      "input-success": allOk,
      "text-success": allOk,
    }}
    type="number"
    step="1"
    min="0"
    bind:value
  />
  <p
    class={{
      label: true,
      "text-wrap": true,
      "text-success": allOk,
    }}
  >
    The number of items to look up for. If value is 0, source default will be
    used.
  </p>
</fieldset>
