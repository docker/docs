using System;
using System.IO;

namespace Replace
{
    class Program
    {
        static void Main(string[] args)
        {
            try
            {
                if (args.Length < 3)
                {
                    Console.WriteLine("Replace.exe [-U] file oldText newText");
                    Environment.Exit(1);
                }

                var replaceText = args[args.Length - 1];
                if (args[0] == "-U")
                    replaceText = replaceText.ToUpper();

                var filePath = args[args.Length - 3];
                if (new FileInfo(filePath).Exists == false)
                    filePath = Path.Combine(Directory.GetCurrentDirectory(), filePath);
                if (new FileInfo(filePath).Exists == false)
                {
                    Console.WriteLine($"Cannot find input file: {args[args.Length - 1]}");
                    Environment.Exit(2);
                }

                var inputText = File.ReadAllText(filePath);
                var outputText = inputText.Replace(args[args.Length - 2], replaceText);
                File.WriteAllText(filePath, outputText);
            }
            catch(Exception ex)
            {
                Console.WriteLine($"Error: {ex.Message}");
            }
        }

    }
}
