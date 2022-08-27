import {Data, DataType} from "../model/Data";

class Client {
    private static _median: Data;
    private static _max: Data;
    private static _min: Data;
    private static _today: Array<Data> = new Array<Data>(10);
    private static _type: DataType = DataType.CARBON_MONOXIDE;
    private static client: Client = new Client();

    private constructor() {
        Client.build();
    }

    private static build() {
        Client._median = Data.build(Client._type);
        Client._max = Data.build(Client._type);
        Client._min = Data.build(Client._type);

        Client._today = new Array<Data>();
        for (let i = 0; i < 10; i++) {
            Client._today.push(Data.build(Client._type));
        }

        Client._today.sort((a: Data, b: Data) => {
            return +b.timestamp - +a.timestamp;
        });
    }

    public static getInstance() {
        return this.client;
    }

    public get latest() {
        return Client._today[0];
    }

    public get median() {
        return Client._median;
    }

    public get max() {
        return Client._max;
    }

    public get min() {
        return Client._min;
    }

    public get today() {
        return Client._today;
    }

    public get type() {
        return Client._type;
    }

    public set type(t: DataType) {
        Client._type = t;
        Client.build();
    }
}

export default Client;