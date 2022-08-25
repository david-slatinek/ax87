class Data {
    dataType: string;
    value: number;
    private readonly timestamp: Date;
    category: number;
    private readonly symbol: string

    constructor(dataType: string, value: number, timestamp: Date, category: number) {
        this.dataType = dataType;
        this.value = value;
        this.timestamp = timestamp;
        this.category = category;
        this.symbol = Data.getSymbol(this.dataType);
    }

    public static getSymbol(dataType: string) {
        switch (dataType) {
            case "CARBON_MONOXIDE":
                return "ppm"
            case "AIR_QUALITY":
                return "aqi"
            default:
                return ""
        }
    }

    public static build() {
        return new Data("CARBON_MONOXIDE", 42, new Date(), 2);
    }

    public getTimestamp() {
        return this.timestamp.toLocaleString();
    }

    public getValueWithSymbol() {
        return this.value + this.symbol;
    }
}

export default Data;