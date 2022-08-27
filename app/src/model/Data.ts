export enum DataType {
    CARBON_MONOXIDE = "CARBON_MONOXIDE",
    AIR_QUALITY = "AIR_QUALITY",
    RAINDROPS = "RAINDROPS",
    SOIL_MOISTURE = "SOIL_MOISTURE"
}

export class Data {
    private readonly dataType: DataType;
    private readonly value: number;
    private readonly _timestamp: Date;
    private readonly category: number;
    private readonly symbol: string
    private readonly categorySymbol: string;

    constructor(dataType: DataType, value: number, timestamp: Date) {
        this.dataType = dataType;
        this.value = +value.toFixed(2);
        this._timestamp = timestamp;
        this.category = +Data.getCategory(dataType, value).toFixed(2);
        this.symbol = Data.getSymbol(dataType);
        this.categorySymbol = Data.getCategorySymbol(dataType);
    }

    public static getRandomDataType() {
        return Object.values(DataType)[Math.floor(Math.random() * Object.keys(DataType).length)];
    }

    private static getRandomValue(dataType: DataType) {
        switch (dataType) {
            case DataType.CARBON_MONOXIDE:
                return Math.floor(Math.random() * (900 - 10) + 10);
            case DataType.AIR_QUALITY:
                return Math.floor(Math.random() * (400 - 20) + 20);
            case DataType.RAINDROPS:
                return Math.floor(Math.random() * 1024);
            case DataType.SOIL_MOISTURE:
                return Math.floor(Math.random() * (489 - 238) + 20);
        }
    }

    private static getRandomDate() {
        const start = new Date();
        start.setDate(start.getDate() - 1);
        return new Date(start.getTime() + Math.random() * (new Date().getTime() - start.getTime()));
    }

    private static getCarbonMonoxideCategory(value: number) {
        if (value <= 30) {
            return 1
        } else if (value > 30 && value <= 70) {
            return 2
        } else if (value > 70 && value <= 150) {
            return 3
        } else if (value > 150 && value <= 200) {
            return 4
        } else if (value > 200 && value <= 400) {
            return 5
        } else if (value > 400 && value <= 800) {
            return 6
        }
        return 7
    }

    private static getAirCategory(value: number) {
        if (value <= 50) {
            return 1
        } else if (value > 50 && value <= 100) {
            return 2
        } else if (value > 100 && value <= 150) {
            return 3
        } else if (value > 150 && value <= 200) {
            return 4
        } else if (value > 200 && value <= 300) {
            return 5
        }
        return 6
    }

    private static mapValue(x: number, inMin: number, inMax: number, outMin: number, outMax: number) {
        return (Math.round(x - inMin) * (outMax - outMin) / (inMax - inMin) + outMin) + 0.5
    }

    private static getCategory(dataType: DataType, value: number) {
        switch (dataType) {
            case DataType.CARBON_MONOXIDE:
                return Data.getCarbonMonoxideCategory(value);
            case DataType.AIR_QUALITY:
                return +Data.getAirCategory(value).toFixed(2);
            case DataType.RAINDROPS:
                return Data.mapValue(value, 0, 1024, 1, 4);
            case DataType.SOIL_MOISTURE:
                return Data.mapValue(value, 489, 238, 0, 100);
        }
        return -1;
    }

    private static getSymbol(dataType: DataType) {
        switch (dataType) {
            case DataType.CARBON_MONOXIDE:
                return "ppm"
            case DataType.AIR_QUALITY:
                return "aqi"
            default:
                return ""
        }
    }

    public static getCategorySymbol(dataType: DataType) {
        switch (dataType) {
            case DataType.RAINDROPS:
            case DataType.SOIL_MOISTURE:
                return "%";
            default:
                return "";
        }
    }

    public static build(dataType: DataType) {
        return new Data(dataType, Data.getRandomValue(dataType), Data.getRandomDate());
    }

    public getTimestamp() {
        return this._timestamp.toLocaleString();
    }

    public get timestamp() {
        return this._timestamp;
    }

    public getValueWithSymbol() {
        return this.value + this.symbol;
    }

    public getCategoryWithSymbol() {
        return this.category + this.categorySymbol;
    }
}