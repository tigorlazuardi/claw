export type SourceServiceClient = Awaited<
  ReturnType<typeof getSourceServiceClient>
>;
type Option = {
  signal?: AbortSignal;
  credentials?: RequestCredentials;
};

export async function getSourceServiceClient(options?: Option) {
  const { createClient } = await import("@connectrpc/connect");
  const { createConnectTransport } = await import("@connectrpc/connect-web");
  const { SourceService } = await import("../gen/claw/v1/source_service_pb");
  const transport = createConnectTransport({
    baseUrl: import.meta.env.BASE_URL,
    fetch: (input, init) => {
      return fetch(input, {
        ...init,
        ...options,
        credentials: options?.credentials || "include",
      });
    },
  });
  return createClient(SourceService, transport);
}
