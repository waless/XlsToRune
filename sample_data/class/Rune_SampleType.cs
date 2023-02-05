using System;
using UnityEngine;
using UnityEngine.AddressableAssets;
using UnityEngine.ResourceManagement.AsyncOperations;
using RuneImporter;

namespace RuneImporter
{
    public static partial class RuneLoader
    {
        public static AsyncOperationHandle Rune_SampleType_LoadInstanceAsync()
        {
            return Rune_SampleType.LoadInstanceAsync();
        }
    }
}

public class Rune_SampleType : RuneScriptableObject
{
    public static Rune_SampleType instance { get; private set; }

    [SerializeField]
    public Value[] ValueList = new Value[4];

    [Serializable]
    public struct Value
    {
        public string name;
        public int number;
        public float position;
    }

    public static AsyncOperationHandle LoadInstanceAsync() {
        var path = Config.ScriptableObjectDirectory + "SampleType.asset";
        var handle = Config.OnLoad(path);
        handle.Completed += (handle) => { instance = handle.Result as Rune_SampleType; };

        return handle;
    }
}
