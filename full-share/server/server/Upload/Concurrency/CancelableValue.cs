namespace server.Upload.Concurrency {
    public struct CancelableValue<T> {
        public T Value { get; set; }
        public bool Cancel { get; set; }

        public static CancelableValue<T> CreateWithValue(T value) {
            return new CancelableValue<T> {
                Value = value,
                Cancel = false
            };
        }

        public static CancelableValue<T> CreateCanceledValue() {
            return new CancelableValue<T> {Cancel = true};
        }
    }
}