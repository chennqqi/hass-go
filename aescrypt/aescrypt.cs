using System;
using System.IO;
using System.Text;
using System.Collections.Generic;
using System.Reflection;
using System.ComponentModel;
using System.Text.RegularExpressions;
using System.Security.Cryptography;

namespace RC4
{
    class crypto
    {
        static int Main(string[] args)
        {
            // decrypt KEY  (will decrypt all .aes files)
            // encrypt KEY  (will encrypt all .txt files)
            // If extension of file ends with .txt => Encrypt
            // If extension of file ends with .aes => Decrypt
            List<string> filenames = new List<string>();
            if (args.Length == 2 && args[0] == "e")
            {
                // Search all .txt files
                FileSelector fs = new FileSelector("name = *.TOML OR name = *.txt OR name = *.jpg OR name = *.tif OR name = *.pdf OR name = *.docx OR name = *.doc OR name = *.xlsx");
                foreach (string filename in fs.SelectFiles(Environment.CurrentDirectory, true))
                {
                    filenames.Add(filename);
                }
            }
            else if (args.Length == 2 && args[0] == "d")
            {
                // Search all .rc4 files
                FileSelector fs = new FileSelector("name = *.aes");
                foreach (string filename in fs.SelectFiles(Environment.CurrentDirectory, true))
                {
                    filenames.Add(filename);
                }
            }
            else
            {
                Console.WriteLine("aescrypt CMD{e|d} KEY");
                Environment.ExitCode = -1;
            }

            if (filenames.Count > 0)
            {
                string key = args[1];

                foreach (string filename in filenames)
                {
                    if (Path.GetExtension(filename) == ".aes")
                    {
                        // Decrypt
                        string aes_filename = filename;
                        string main_filename = aes_filename.Substring(0, aes_filename.Length - 4);
                        if (!File.Exists(main_filename) || (File.GetLastWriteTime(aes_filename) != File.GetLastWriteTime(main_filename)))
                        {
                            Console.WriteLine("Decrypted \"{0}\"", main_filename);
                            try
                            {
                                using (Stream inputstream = File.OpenRead(aes_filename))
                                using (Stream outputstream = File.OpenWrite(main_filename))
                                    SharpAESCrypt.SharpAESCrypt.Decrypt(key, inputstream, outputstream);
                            }
                            catch (Exception ex)
                            {
                                Console.WriteLine(string.Format(SharpAESCrypt.Strings.CommandlineError, ex.ToString()));
                            }
                            File.SetLastWriteTime(main_filename, File.GetLastWriteTime(aes_filename));
                        }
                        else
                        {
                            Console.WriteLine("Unchanged \"{0}\"", main_filename);
                        }
                    }
                    else
                    {
                        // Encrypt
                        string main_filename = filename;
                        string aes_filename = filename + ".aes";
                        if (!File.Exists(aes_filename) || (File.GetLastWriteTime(aes_filename) != File.GetLastWriteTime(main_filename)))
                        {
                            Console.WriteLine("Encrypted \"{0}\"", main_filename);
                            try
                            {
                                using (Stream inputstream = File.OpenRead(filename))
                                using (Stream outputstream = File.OpenWrite(filename + ".aes"))
                                    SharpAESCrypt.SharpAESCrypt.Encrypt(key, inputstream, outputstream);
                            }
                            catch (Exception ex)
                            {
                                Console.WriteLine(string.Format(SharpAESCrypt.Strings.CommandlineError, ex.ToString()));
                            }
                            File.SetLastWriteTime(aes_filename, File.GetLastWriteTime(main_filename));
                        }
                        else
                        {
                            Console.WriteLine("Unchanged \"{0}\"", main_filename);
                        }
                    }

                }
            }

            return Environment.ExitCode;
        }

    }




    /// <summary>
    /// Enumerates the options for a logical conjunction. This enum is intended for use
    /// internally by the FileSelector class.
    /// </summary>
    internal enum LogicalConjunction
    {
        NONE,
        AND,
        OR,
        XOR,
    }

    internal enum WhichTime
    {
        atime,
        mtime,
        ctime,
    }


    internal enum ComparisonOperator
    {
        [Description(">")]
        GreaterThan,
        [Description(">=")]
        GreaterThanOrEqualTo,
        [Description("<")]
        LesserThan,
        [Description("<=")]
        LesserThanOrEqualTo,
        [Description("=")]
        EqualTo,
        [Description("!=")]
        NotEqualTo
    }


    internal abstract partial class SelectionCriterion
    {
        internal abstract bool Evaluate(string filename);
    }


    internal partial class SizeCriterion : SelectionCriterion
    {
        internal ComparisonOperator Operator;
        internal Int64 Size;

        public override String ToString()
        {
            StringBuilder sb = new StringBuilder();
            sb.Append("size ").Append(EnumUtil.GetDescription(Operator)).Append(" ").Append(Size.ToString());
            return sb.ToString();
        }

        internal override bool Evaluate(string filename)
        {
            System.IO.FileInfo fi = new System.IO.FileInfo(filename);
            return _Evaluate(fi.Length);
        }

        private bool _Evaluate(Int64 Length)
        {
            bool result = false;
            switch (Operator)
            {
                case ComparisonOperator.GreaterThanOrEqualTo:
                    result = Length >= Size;
                    break;
                case ComparisonOperator.GreaterThan:
                    result = Length > Size;
                    break;
                case ComparisonOperator.LesserThanOrEqualTo:
                    result = Length <= Size;
                    break;
                case ComparisonOperator.LesserThan:
                    result = Length < Size;
                    break;
                case ComparisonOperator.EqualTo:
                    result = Length == Size;
                    break;
                case ComparisonOperator.NotEqualTo:
                    result = Length != Size;
                    break;
                default:
                    throw new ArgumentException("Operator");
            }
            return result;
        }

    }



    internal partial class TimeCriterion : SelectionCriterion
    {
        internal ComparisonOperator Operator;
        internal WhichTime Which;
        internal DateTime Time;

        public override String ToString()
        {
            StringBuilder sb = new StringBuilder();
            sb.Append(Which.ToString()).Append(" ").Append(EnumUtil.GetDescription(Operator)).Append(" ").Append(Time.ToString("yyyy-MM-dd-HH:mm:ss"));
            return sb.ToString();
        }

        internal override bool Evaluate(string filename)
        {
            System.IO.FileInfo fi = new System.IO.FileInfo(filename);
            DateTime x;
            switch (Which)
            {
                case WhichTime.atime:
                    x = System.IO.File.GetLastAccessTime(filename);
                    break;
                case WhichTime.mtime:
                    x = System.IO.File.GetLastWriteTime(filename);
                    break;
                case WhichTime.ctime:
                    x = System.IO.File.GetCreationTime(filename);
                    break;
                default:
                    throw new ArgumentException("Operator");
            }
            return _Evaluate(x);
        }


        private bool _Evaluate(DateTime x)
        {

            bool result = false;
            switch (Operator)
            {
                case ComparisonOperator.GreaterThanOrEqualTo:
                    result = (x >= Time);
                    break;
                case ComparisonOperator.GreaterThan:
                    result = (x > Time);
                    break;
                case ComparisonOperator.LesserThanOrEqualTo:
                    result = (x <= Time);
                    break;
                case ComparisonOperator.LesserThan:
                    result = (x < Time);
                    break;
                case ComparisonOperator.EqualTo:
                    result = (x == Time);
                    break;
                case ComparisonOperator.NotEqualTo:
                    result = (x != Time);
                    break;
                default:
                    throw new ArgumentException("Operator");
            }

            //Console.WriteLine("TimeCriterion[{2}]({0})= {1}", filename, result, Which.ToString());
            return result;
        }
    }

    internal partial class NameCriterion : SelectionCriterion
    {
        private Regex _re;
        private String _regexString;
        internal ComparisonOperator Operator;
        private string _MatchingFileSpec;
        internal virtual string MatchingFileSpec
        {
            set
            {
                _MatchingFileSpec = value;
                _regexString = "^" +
                Regex.Escape(value)
                .Replace(@"\*\.\*", @"([^\.]+|.*\.[^\\\.]*)")
                .Replace(@"\.\*", @"\.[^\\\.]*")
                .Replace(@"\*", @".*")
                .Replace(@"\?", @"[^\\\.]")
                + "$";

                // neither of these is correct
                //if (!_regexString.StartsWith(@"\\")) _regexString = @"\\" + _regexString;
                //if (_regexString.IndexOf("\\") == -1)  _regexString = @"\\" + _regexString;
                _re = new Regex(_regexString, RegexOptions.IgnoreCase);
            }
        }


        public override String ToString()
        {
            StringBuilder sb = new StringBuilder();
            sb.Append("name = ").Append(_MatchingFileSpec);
            return sb.ToString();
        }


        internal override bool Evaluate(string filename)
        {
            return _Evaluate(filename);
        }

        private bool _Evaluate(string fullpath)
        {
            // No slash in the pattern implicitly means recurse, which means compare to
            // filename only, not full path.
            String f = (_MatchingFileSpec.IndexOf('\\') == -1)
                ? System.IO.Path.GetFileName(fullpath)
                : fullpath; // compare to fullpath

            bool result = _re.IsMatch(f);
            if (Operator != ComparisonOperator.EqualTo)
                result = !result;
            return result;
        }
    }



    internal partial class AttributesCriterion : SelectionCriterion
    {
        private FileAttributes _Attributes;
        internal ComparisonOperator Operator;
        internal string AttributeString
        {
            get
            {
                string result = "";
                if ((_Attributes & FileAttributes.Hidden) == FileAttributes.Hidden)
                    result += "H";
                if ((_Attributes & FileAttributes.System) == FileAttributes.System)
                    result += "S";
                if ((_Attributes & FileAttributes.ReadOnly) == FileAttributes.ReadOnly)
                    result += "R";
                if ((_Attributes & FileAttributes.Archive) == FileAttributes.Archive)
                    result += "A";
                if ((_Attributes & FileAttributes.NotContentIndexed) == FileAttributes.NotContentIndexed)
                    result += "I";
                return result;
            }

            set
            {
                _Attributes = FileAttributes.Normal;
                foreach (char c in value.ToUpper())
                {
                    switch (c)
                    {
                        case 'H':
                            if ((_Attributes & FileAttributes.Hidden) == FileAttributes.Hidden)
                                throw new ArgumentException(String.Format("Repeated flag. ({0})", c), "value");
                            _Attributes |= FileAttributes.Hidden;
                            break;

                        case 'R':
                            if ((_Attributes & FileAttributes.ReadOnly) == FileAttributes.ReadOnly)
                                throw new ArgumentException(String.Format("Repeated flag. ({0})", c), "value");
                            _Attributes |= FileAttributes.ReadOnly;
                            break;

                        case 'S':
                            if ((_Attributes & FileAttributes.System) == FileAttributes.System)
                                throw new ArgumentException(String.Format("Repeated flag. ({0})", c), "value");
                            _Attributes |= FileAttributes.System;
                            break;

                        case 'A':
                            if ((_Attributes & FileAttributes.Archive) == FileAttributes.Archive)
                                throw new ArgumentException(String.Format("Repeated flag. ({0})", c), "value");
                            _Attributes |= FileAttributes.Archive;
                            break;

                        case 'I':
                            if ((_Attributes & FileAttributes.NotContentIndexed) == FileAttributes.NotContentIndexed)
                                throw new ArgumentException(String.Format("Repeated flag. ({0})", c), "value");
                            _Attributes |= FileAttributes.NotContentIndexed;
                            break;
                        default:
                            throw new ArgumentException(value);
                    }
                }
            }
        }


        public override String ToString()
        {
            StringBuilder sb = new StringBuilder();
            sb.Append("attributes ").Append(EnumUtil.GetDescription(Operator)).Append(" ").Append(AttributeString);
            return sb.ToString();
        }

        private bool _EvaluateOne(FileAttributes fileAttrs, FileAttributes criterionAttrs)
        {
            bool result = false;
            if ((_Attributes & criterionAttrs) == criterionAttrs)
                result = ((fileAttrs & criterionAttrs) == criterionAttrs);
            else
                result = true;
            return result;
        }



        internal override bool Evaluate(string filename)
        {
#if NETCF
		FileAttributes fileAttrs = NetCfFile.GetAttributes(filename);
#else
            FileAttributes fileAttrs = System.IO.File.GetAttributes(filename);
#endif

            return _Evaluate(fileAttrs);
        }

        private bool _Evaluate(FileAttributes fileAttrs)
        {
            //Console.WriteLine("fileattrs[{0}]={1}", filename, fileAttrs.ToString());

            bool result = _EvaluateOne(fileAttrs, FileAttributes.Hidden);
            if (result)
                result = _EvaluateOne(fileAttrs, FileAttributes.System);
            if (result)
                result = _EvaluateOne(fileAttrs, FileAttributes.ReadOnly);
            if (result)
                result = _EvaluateOne(fileAttrs, FileAttributes.Archive);

            if (Operator != ComparisonOperator.EqualTo)
                result = !result;

            //Console.WriteLine("AttributesCriterion[{2}]({0})= {1}", filename, result, AttributeString);

            return result;
        }
    }



    internal partial class CompoundCriterion : SelectionCriterion
    {
        internal LogicalConjunction Conjunction;
        internal SelectionCriterion Left;

        private SelectionCriterion _Right;
        internal SelectionCriterion Right
        {
            get { return _Right; }
            set
            {
                _Right = value;
                if (value == null)
                    Conjunction = LogicalConjunction.NONE;
                else if (Conjunction == LogicalConjunction.NONE)
                    Conjunction = LogicalConjunction.AND;
            }
        }


        internal override bool Evaluate(string filename)
        {
            bool result = Left.Evaluate(filename);
            switch (Conjunction)
            {
                case LogicalConjunction.AND:
                    if (result)
                        result = Right.Evaluate(filename);
                    break;
                case LogicalConjunction.OR:
                    if (!result)
                        result = Right.Evaluate(filename);
                    break;
                case LogicalConjunction.XOR:
                    result ^= Right.Evaluate(filename);
                    break;
                default:
                    throw new ArgumentException("Conjunction");
            }
            return result;
        }


        public override String ToString()
        {
            StringBuilder sb = new StringBuilder();
            sb.Append("(")
            .Append((Left != null) ? Left.ToString() : "null")
            .Append(" ")
            .Append(Conjunction.ToString())
            .Append(" ")
            .Append((Right != null) ? Right.ToString() : "null")
            .Append(")");
            return sb.ToString();
        }
    }



    /// <summary>
    /// FileSelector encapsulates logic that selects files from a source based on a set
    /// of criteria.
    /// </summary>
    /// <remarks>
    ///
    /// <para>
    /// But, some applications may wish to use the FileSelector class directly, to select
    /// files from disk volumes based on a set of criteria, without creating or querying Zip
    /// archives.  The file selection criteria include: a pattern to match the filename; the
    /// last modified, created, or last accessed time of the file; the size of the file; and
    /// the attributes of the file.
    /// </para>
    /// </remarks>
    public partial class FileSelector
    {
        internal SelectionCriterion _Criterion;

        /// <summary>
        /// The default constructor.
        /// </summary>
        /// <remarks>
        /// Typically, applications won't use this constructor.  Instead they'll call the
        /// constructor that accepts a selectionCriteria string.  If you use this constructor,
        /// you'll want to set the SelectionCriteria property on the instance before calling
        /// SelectFiles().
        /// </remarks>
        protected FileSelector() { }


        /// <summary>
        /// Constructor that allows the caller to specify file selection criteria.
        /// </summary>
        ///
        /// <remarks>
        /// <para>
        /// This constructor allows the caller to specify a set of criteria for selection of files.
        /// </para>
        ///
        /// <para>
        /// See <see cref="FileSelector.SelectionCriteria"/> for a description of the syntax of
        /// the selectionCriteria string.
        /// </para>
        /// </remarks>
        ///
        /// <param name="selectionCriteria">The criteria for file selection.</param>
        public FileSelector(String selectionCriteria)
        {
            if (!String.IsNullOrEmpty(selectionCriteria))
                _Criterion = _ParseCriterion(selectionCriteria);
        }



        /// <summary>
        /// The string specifying which files to include when retrieving.
        /// </summary>
        /// <remarks>
        ///
        /// <para>
        /// Specify the criteria in statements of 3 elements: a noun, an operator, and a value.
        /// Consider the string "name != *.doc" .  The noun is "name".  The operator is "!=",
        /// implying "Not Equal".  The value is "*.doc".  That criterion, in English, says "all
        /// files with a name that does not end in the .doc extension."
        /// </para>
        ///
        /// <para>
        /// Supported nouns include "name" for the filename; "atime", "mtime", and "ctime" for
        /// last access time, last modfied time, and created time of the file, respectively;
        /// "attributes" for the file attributes; and "size" for the file length (uncompressed).
        /// The "attributes" and "name" nouns both support = and != as operators.  The "size",
        /// "atime", "mtime", and "ctime" nouns support = and !=, and &gt;, &gt;=, &lt;, &lt;=
        /// as well.
        /// </para>
        ///
        /// <para>
        /// Specify values for the file attributes as a string with one or more of the
        /// characters H,R,S,A,I in any order, implying Hidden, ReadOnly, System, Archive,
        /// and NotContextIndexed,
        /// respectively.  To specify a time, use YYYY-MM-DD-HH:mm:ss as the format.  If you
        /// omit the HH:mm:ss portion, it is assumed to be 00:00:00 (midnight). The value for a
        /// size criterion is expressed in integer quantities of bytes, kilobytes (use k or kb
        /// after the number), megabytes (m or mb), or gigabytes (g or gb).  The value for a
        /// name is a pattern to match against the filename, potentially including wildcards.
        /// The pattern follows CMD.exe glob rules: * implies one or more of any character,
        /// while ? implies one character.  If the name pattern contains any slashes, it is
        /// matched to the entire filename, including the path; otherwise, it is matched
        /// against only the filename without the path.  This means a pattern of "*\*.*" matches
        /// all files one directory level deep, while a pattern of "*.*" matches all files in
        /// all directories.
        /// </para>
        ///
        /// <para>
        /// To specify a name pattern that includes spaces, use single quotes around the pattern.
        /// A pattern of "'* *.*'" will match all files that have spaces in the filename.  The full
        /// criteria string for that would be "name = '* *.*'" .
        /// </para>
        ///
        /// <para>
        /// Some examples: a string like "attributes != H" retrieves all entries whose
        /// attributes do not include the Hidden bit.  A string like "mtime > 2009-01-01"
        /// retrieves all entries with a last modified time after January 1st, 2009.  For
        /// example "size &gt; 2gb" retrieves all entries whose uncompressed size is greater
        /// than 2gb.
        /// </para>
        ///
        /// <para>
        /// You can combine criteria with the conjunctions AND, OR, and XOR. Using a string like
        /// "name = *.txt AND size &gt;= 100k" for the selectionCriteria retrieves entries whose
        /// names end in .txt, and whose uncompressed size is greater than or equal to 100
        /// kilobytes.
        /// </para>
        ///
        /// <para>
        /// For more complex combinations of criteria, you can use parenthesis to group clauses
        /// in the boolean logic.  Absent parenthesis, the precedence of the criterion atoms is
        /// determined by order of appearance.  Unlike the C# language, the AND conjunction does
        /// not take precendence over the logical OR.  This is important only in strings that
        /// contain 3 or more criterion atoms.  In other words, "name = *.txt and size &gt; 1000
        /// or attributes = H" implies "((name = *.txt AND size &gt; 1000) OR attributes = H)"
        /// while "attributes = H OR name = *.txt and size &gt; 1000" evaluates to "((attributes
        /// = H OR name = *.txt) AND size &gt; 1000)".  When in doubt, use parenthesis.
        /// </para>
        ///
        /// <para>
        /// Using time properties requires some extra care. If you want to retrieve all entries
        /// that were last updated on 2009 February 14, specify "mtime &gt;= 2009-02-14 AND
        /// mtime &lt; 2009-02-15".  Read this to say: all files updated after 12:00am on
        /// February 14th, until 12:00am on February 15th.  You can use the same bracketing
        /// approach to specify any time period - a year, a month, a week, and so on.
        /// </para>
        ///
        /// <para>
        /// The syntax allows one special case: if you provide a string with no spaces, it is treated as
        /// a pattern to match for the filename.  Therefore a string like "*.xls" will be equivalent to
        /// specifying "name = *.xls".
        /// </para>
        ///
        /// <para>
        /// There is no logic in this class that insures that the inclusion criteria
        /// are internally consistent.  For example, it's possible to specify criteria that
        /// says the file must have a size of less than 100 bytes, as well as a size that
        /// is greater than 1000 bytes.  Obviously no file will ever satisfy such criteria,
        /// but this class does not check and find such inconsistencies.
        /// </para>
        ///
        /// </remarks>
        ///
        /// <exception cref="System.Exception">
        /// Thrown in the setter if the value has an invalid syntax.
        /// </exception>
        public String SelectionCriteria
        {
            get
            {
                if (_Criterion == null) return null;
                return _Criterion.ToString();
            }
            set
            {
                if (value == null) _Criterion = null;
                else if (value.Trim() == "") _Criterion = null;
                else
                    _Criterion = _ParseCriterion(value);
            }
        }


        private enum ParseState
        {
            Start,
            OpenParen,
            CriterionDone,
            ConjunctionPending,
            Whitespace,
        }



        private SelectionCriterion _ParseCriterion(String s)
        {
            if (s == null) return null;

            // shorthand for filename glob
            if (s.IndexOf(" ") == -1)
                s = "name = " + s;

            // inject spaces after open paren and before close paren
            string[] prPairs = { @"\((\S)", "( $1", @"(\S)\)", "$1 )", };
            for (int i = 0; i + 1 < prPairs.Length; i += 2)
            {
                Regex rgx = new Regex(prPairs[i]);
                s = rgx.Replace(s, prPairs[i + 1]);
            }

            // split the expression into tokens
            string[] tokens = s.Trim().Split(' ', '\t');

            if (tokens.Length < 3) throw new ArgumentException(s);

            SelectionCriterion current = null;

            LogicalConjunction pendingConjunction = LogicalConjunction.NONE;

            ParseState state;
            var stateStack = new System.Collections.Generic.Stack<ParseState>();
            var critStack = new System.Collections.Generic.Stack<SelectionCriterion>();
            stateStack.Push(ParseState.Start);

            for (int i = 0; i < tokens.Length; i++)
            {
                switch (tokens[i].ToLower())
                {
                    case "and":
                    case "xor":
                    case "or":
                        state = stateStack.Peek();
                        if (state != ParseState.CriterionDone)
                            throw new ArgumentException(String.Join(" ", tokens, i, tokens.Length - i));

                        if (tokens.Length <= i + 3)
                            throw new ArgumentException(String.Join(" ", tokens, i, tokens.Length - i));

                        pendingConjunction = (LogicalConjunction)Enum.Parse(typeof(LogicalConjunction), tokens[i].ToUpper());
                        current = new CompoundCriterion { Left = current, Right = null, Conjunction = pendingConjunction };
                        stateStack.Push(state);
                        stateStack.Push(ParseState.ConjunctionPending);
                        critStack.Push(current);
                        break;

                    case "(":
                        state = stateStack.Peek();
                        if (state != ParseState.Start && state != ParseState.ConjunctionPending && state != ParseState.OpenParen)
                            throw new ArgumentException(String.Join(" ", tokens, i, tokens.Length - i));

                        if (tokens.Length <= i + 4)
                            throw new ArgumentException(String.Join(" ", tokens, i, tokens.Length - i));

                        stateStack.Push(ParseState.OpenParen);
                        break;

                    case ")":
                        state = stateStack.Pop();
                        if (stateStack.Peek() != ParseState.OpenParen)
                            throw new ArgumentException(String.Join(" ", tokens, i, tokens.Length - i));

                        stateStack.Pop();
                        stateStack.Push(ParseState.CriterionDone);
                        break;

                    case "atime":
                    case "ctime":
                    case "mtime":
                        if (tokens.Length <= i + 2)
                            throw new ArgumentException(String.Join(" ", tokens, i, tokens.Length - i));

                        DateTime t;
                        try
                        {
                            t = DateTime.ParseExact(tokens[i + 2], "yyyy-MM-dd-HH:mm:ss", null);
                        }
                        catch
                        {
                            t = DateTime.ParseExact(tokens[i + 2], "yyyy-MM-dd", null);
                        }
                        current = new TimeCriterion
                        {
                            Which = (WhichTime)Enum.Parse(typeof(WhichTime), tokens[i]),
                            Operator = (ComparisonOperator)EnumUtil.Parse(typeof(ComparisonOperator), tokens[i + 1]),
                            Time = t
                        };
                        i += 2;
                        stateStack.Push(ParseState.CriterionDone);
                        break;


                    case "length":
                    case "size":
                        if (tokens.Length <= i + 2)
                            throw new ArgumentException(String.Join(" ", tokens, i, tokens.Length - i));

                        Int64 sz = 0;
                        string v = tokens[i + 2];
                        if (v.ToUpper().EndsWith("K"))
                            sz = Int64.Parse(v.Substring(0, v.Length - 1)) * 1024;
                        else if (v.ToUpper().EndsWith("KB"))
                            sz = Int64.Parse(v.Substring(0, v.Length - 2)) * 1024;
                        else if (v.ToUpper().EndsWith("M"))
                            sz = Int64.Parse(v.Substring(0, v.Length - 1)) * 1024 * 1024;
                        else if (v.ToUpper().EndsWith("MB"))
                            sz = Int64.Parse(v.Substring(0, v.Length - 2)) * 1024 * 1024;
                        else if (v.ToUpper().EndsWith("G"))
                            sz = Int64.Parse(v.Substring(0, v.Length - 1)) * 1024 * 1024 * 1024;
                        else if (v.ToUpper().EndsWith("GB"))
                            sz = Int64.Parse(v.Substring(0, v.Length - 2)) * 1024 * 1024 * 1024;
                        else sz = Int64.Parse(tokens[i + 2]);

                        current = new SizeCriterion
                        {
                            Size = sz,
                            Operator = (ComparisonOperator)EnumUtil.Parse(typeof(ComparisonOperator), tokens[i + 1])
                        };
                        i += 2;
                        stateStack.Push(ParseState.CriterionDone);
                        break;

                    case "filename":
                    case "name":
                        {
                            if (tokens.Length <= i + 2)
                                throw new ArgumentException(String.Join(" ", tokens, i, tokens.Length - i));

                            ComparisonOperator c =
                                (ComparisonOperator)EnumUtil.Parse(typeof(ComparisonOperator), tokens[i + 1]);

                            if (c != ComparisonOperator.NotEqualTo && c != ComparisonOperator.EqualTo)
                                throw new ArgumentException(String.Join(" ", tokens, i, tokens.Length - i));

                            string m = tokens[i + 2];
                            // handle single-quoted filespecs (used to include spaces in filename patterns)
                            if (m.StartsWith("'"))
                            {
                                int ix = i;
                                if (!m.EndsWith("'"))
                                {
                                    do
                                    {
                                        i++;
                                        if (tokens.Length <= i + 2)
                                            throw new ArgumentException(String.Join(" ", tokens, ix, tokens.Length - ix));
                                        m += " " + tokens[i + 2];
                                    } while (!tokens[i + 2].EndsWith("'"));
                                }
                                // trim off leading and trailing single quotes
                                m = m.Substring(1, m.Length - 2);
                            }

                            current = new NameCriterion
                            {
                                MatchingFileSpec = m,
                                Operator = c
                            };
                            i += 2;
                            stateStack.Push(ParseState.CriterionDone);
                        }
                        break;

                    case "attributes":
                        {
                            if (tokens.Length <= i + 2)
                                throw new ArgumentException(String.Join(" ", tokens, i, tokens.Length - i));

                            ComparisonOperator c =
                                (ComparisonOperator)EnumUtil.Parse(typeof(ComparisonOperator), tokens[i + 1]);

                            if (c != ComparisonOperator.NotEqualTo && c != ComparisonOperator.EqualTo)
                                throw new ArgumentException(String.Join(" ", tokens, i, tokens.Length - i));

                            current = new AttributesCriterion
                            {
                                AttributeString = tokens[i + 2],
                                Operator = c
                            };
                            i += 2;
                            stateStack.Push(ParseState.CriterionDone);
                        }
                        break;

                    case "":
                        // NOP
                        stateStack.Push(ParseState.Whitespace);
                        break;

                    default:
                        throw new ArgumentException("'" + tokens[i] + "'");
                }

                state = stateStack.Peek();
                if (state == ParseState.CriterionDone)
                {
                    stateStack.Pop();
                    if (stateStack.Peek() == ParseState.ConjunctionPending)
                    {
                        while (stateStack.Peek() == ParseState.ConjunctionPending)
                        {
                            var cc = critStack.Pop() as CompoundCriterion;
                            cc.Right = current;
                            current = cc; // mark the parent as current (walk up the tree)
                            stateStack.Pop();   // the conjunction is no longer pending

                            state = stateStack.Pop();
                            if (state != ParseState.CriterionDone)
                                throw new ArgumentException();
                        }
                    }
                    else stateStack.Push(ParseState.CriterionDone);  // not sure?
                }

                if (state == ParseState.Whitespace)
                    stateStack.Pop();
            }

            return current;
        }


        /// <summary>
        /// Returns a string representation of the FileSelector object.
        /// </summary>
        /// <returns>The string representation of the boolean logic statement of the file
        /// selection criteria for this instance. </returns>
        public override String ToString()
        {
            return _Criterion.ToString();
        }


        private bool Evaluate(string filename)
        {
            bool result = _Criterion.Evaluate(filename);
            return result;
        }


        /// <summary>
        /// Returns the names of the files in the specified directory
        /// that fit the selection criteria specified in the FileSelector.
        /// </summary>
        ///
        /// <remarks>
        /// This is equivalent to calling <see cref="SelectFiles(String, bool)"/>
        /// with recurseDirectories = false.
        /// </remarks>
        ///
        /// <param name="directory">
        /// The name of the directory over which to apply the FileSelector criteria.
        /// </param>
        ///
        /// <returns>
        /// A collection of strings containing fully-qualified pathnames of files
        /// that match the criteria specified in the FileSelector instance.
        /// </returns>
        public System.Collections.Generic.ICollection<String> SelectFiles(String directory)
        {
            return SelectFiles(directory, false);
        }


        /// <summary>
        /// Returns the names of the files in the specified directory that fit the selection
        /// criteria specified in the FileSelector, optionally recursing through subdirectories.
        /// </summary>
        ///
        /// <remarks>
        /// This method applies the file selection criteria contained in the FileSelector to the
        /// files contained in the given directory, and returns the names of files that
        /// conform to the criteria.
        /// </remarks>
        ///
        /// <param name="directory">
        /// The name of the directory over which to apply the FileSelector criteria.
        /// </param>
        ///
        /// <param name="recurseDirectories">
        /// Whether to recurse through subdirectories when applying the file selection criteria.
        /// </param>
        ///
        /// <returns>
        /// An collection of strings containing fully-qualified pathnames of files
        /// that match the criteria specified in the FileSelector instance.
        /// </returns>
        public System.Collections.Generic.ICollection<String> SelectFiles(String directory, bool recurseDirectories)
        {
            if (_Criterion == null)
                throw new ArgumentException("SelectionCriteria has not been set");

            String[] filenames = System.IO.Directory.GetFiles(directory);
            var list = new System.Collections.Generic.List<String>();

            // add the files:
            foreach (String filename in filenames)
            {
                if (Evaluate(filename))
                    list.Add(filename);
            }

            if (recurseDirectories)
            {
                // add the subdirectories:
                String[] dirnames = System.IO.Directory.GetDirectories(directory);
                foreach (String dir in dirnames)
                {
                    DirectoryInfo d = new DirectoryInfo(dir);
                    if (d.Attributes.HasFlag(FileAttributes.Hidden))
                        continue;
                    if (d.Name == ".git" || d.Name == ".hg" || d.Name == ".svn")
                        continue;
                    list.AddRange(this.SelectFiles(dir, recurseDirectories));
                }
            }
            return list;
        }
    }




    /// <summary>
    /// Summary description for EnumUtil.
    /// </summary>
    internal sealed class EnumUtil
    {
        /// <summary>
        /// Returns the value of the DescriptionAttribute if the specified Enum value has one.
        /// If not, returns the ToString() representation of the Enum value.
        /// </summary>
        /// <param name="value">The Enum to get the description for</param>
        /// <returns></returns>
        internal static string GetDescription(System.Enum value)
        {
            FieldInfo fi = value.GetType().GetField(value.ToString());
            var attributes = (DescriptionAttribute[])fi.GetCustomAttributes(typeof(DescriptionAttribute), false);
            if (attributes.Length > 0)
                return attributes[0].Description;
            else
                return value.ToString();
        }

        /// <summary>
        /// Converts the string representation of the name or numeric value of one or more
        /// enumerated constants to an equivalent enumerated object.
        /// Note: use the DescriptionAttribute on enum values to enable this.
        /// </summary>
        /// <param name="enumType">The System.Type of the enumeration.</param>
        /// <param name="stringRepresentation">A string containing the name or value to convert.</param>
        /// <returns></returns>
        internal static object Parse(Type enumType, string stringRepresentation)
        {
            return Parse(enumType, stringRepresentation, false);
        }

        /// <summary>
        /// Converts the string representation of the name or numeric value of one or more
        /// enumerated constants to an equivalent enumerated object.
        /// A parameter specified whether the operation is case-sensitive.
        /// Note: use the DescriptionAttribute on enum values to enable this.
        /// </summary>
        /// <param name="enumType">The System.Type of the enumeration.</param>
        /// <param name="stringRepresentation">A string containing the name or value to convert.</param>
        /// <param name="ignoreCase">Whether the operation is case-sensitive or not.</param>
        /// <returns></returns>
        internal static object Parse(Type enumType, string stringRepresentation, bool ignoreCase)
        {
            if (ignoreCase)
                stringRepresentation = stringRepresentation.ToLower();

            foreach (System.Enum enumVal in System.Enum.GetValues(enumType))
            {
                string description = GetDescription(enumVal);
                if (ignoreCase)
                    description = description.ToLower();
                if (description == stringRepresentation)
                    return enumVal;
            }

            return System.Enum.Parse(enumType, stringRepresentation, ignoreCase);
        }
    }

}


namespace SharpAESCrypt
{
    /// <summary>
    /// Enumerates the possible modes for encryption and decryption
    /// </summary>
    public enum OperationMode
    {
        /// <summary>
        /// Indicates encryption, which means that the stream must be writeable
        /// </summary>
        Encrypt,
        /// <summary>
        /// Indicates decryption, which means that the stream must be readable
        /// </summary>
        Decrypt
    }

    #region Translateable strings
    /// <summary>
    /// Placeholder for translateable strings
    /// </summary>
    public static class Strings
    {
        #region Command line
        /// <summary>
        /// A string displayed when the program is invoked without the correct number of arguments
        /// </summary>
        public static string CommandlineUsage = "SharpAESCrypt e|d <password> [<fromPath>] [<toPath>]" +
            Environment.NewLine +
            Environment.NewLine +
            "If you ommit the fromPath or toPath, stdin/stdout are used insted, e.g.:" +
            Environment.NewLine +
            " SharpAESCrypt e 1234 < file.jpg > file.jpg.aes"
            ;

        /// <summary>
        /// A string displayed when an error occurs while running the commandline program
        /// </summary>
        public static string CommandlineError = "Error: {0}";
        /// <summary>
        /// A string displayed if the mode is neither e nor d
        /// </summary>
        public static string CommandlineUnknownMode = "Invalid operation, must be (e)ncrypt or (d)ecrypt";
        #endregion

        #region Exception messages
        /// <summary>
        /// An exception message that indicates that the hash algorithm is not supported
        /// </summary>
        public static string UnsupportedHashAlgorithmReuse = "The hash algortihm does not support reuse";
        /// <summary>
        /// An exception message that indicates that the hash algorithm is not supported
        /// </summary>
        public static string UnsupportedHashAlgorithmBlocks = "The hash algortihm does not support multiple blocks";
        /// <summary>
        /// An exception message that indicates that the hash algorithm is not supported
        /// </summary>
        public static string UnsupportedHashAlgorithmBlocksize = "Unable to digest {0} bytes, as the hash algorithm only returns {1} bytes";
        /// <summary>
        /// An exception message that indicates that an unexpected end of stream was encountered
        /// </summary>
        public static string UnexpectedEndOfStream = "The stream was exhausted unexpectedly";
        /// <summary>
        /// An exception message that indicates that the stream does not support writing
        /// </summary>
        public static string StreamMustBeWriteAble = "When encrypting, the stream must be writeable";
        /// <summary>
        /// An exception messaget that indicates that the stream does not support reading
        /// </summary>
        public static string StreamMustBeReadAble = "When decrypting, the stream must be readable";
        /// <summary>
        /// An exception message that indicates that the mode is not one of the allowed enumerations
        /// </summary>
        public static string InvalidOperationMode = "Invalid mode supplied";

        /// <summary>
        /// An exception message that indicates that file is not in the correct format
        /// </summary>
        public static string InvalidFileFormat = "Invalid file format";
        /// <summary>
        /// An exception message that indicates that the header marker is invalid
        /// </summary>
        public static string InvalidHeaderMarker = "Invalid header marker";
        /// <summary>
        /// An exception message that indicates that the reserved field is not set to zero
        /// </summary>
        public static string InvalidReservedFieldValue = "Reserved field is not zero";
        /// <summary>
        /// An exception message that indicates that the detected file version is not supported
        /// </summary>
        public static string UnsupportedFileVersion = "Unsuported file version: {0}";
        /// <summary>
        /// An exception message that indicates that an extension had an invalid format
        /// </summary>
        public static string InvalidExtensionData = "Invalid extension data, separator (0x00) not found";
        /// <summary>
        /// An exception message that indicates that the format was accepted, but the password was not verified
        /// </summary>
        public static string InvalidPassword = "Invalid password or corrupted data";
        /// <summary>
        /// An exception message that indicates that the length of the file is incorrect
        /// </summary>
        public static string InvalidFileLength = "File length is invalid";

        /// <summary>
        /// An exception message that indicates that the version is readonly when decrypting
        /// </summary>
        public static string VersionReadonlyForDecryption = "Version is readonly when decrypting";
        /// <summary>
        /// An exception message that indicates that the file version is readonly once encryption has started
        /// </summary>
        public static string VersionReadonly = "Version cannot be changed after encryption has started";
        /// <summary>
        /// An exception message that indicates that the supplied version number is unsupported
        /// </summary>
        public static string VersionUnsupported = "The maximum allowed version is {0}";
        /// <summary>
        /// An exception message that indicates that the stream must support seeking
        /// </summary>
        public static string StreamMustSupportSeeking = "The stream must be seekable writing version 0 files";

        /// <summary>
        /// An exception message that indicates that the requsted operation is unsupported
        /// </summary>
        public static string CannotReadWhileEncrypting = "Cannot read while encrypting";
        /// <summary>
        /// An exception message that indicates that the requsted operation is unsupported
        /// </summary>
        public static string CannotWriteWhileDecrypting = "Cannot read while decrypting";

        /// <summary>
        /// An exception message that indicates that the data has been altered
        /// </summary>
        public static string DataHMACMismatch = "Message has been altered, do not trust content";
        /// <summary>
        /// An exception message that indicates that the data has been altered or the password is invalid
        /// </summary>
        public static string DataHMACMismatch_v0 = "Invalid password or content has been altered";

        /// <summary>
        /// An exception message that indicates that the system is missing a text encoding
        /// </summary>
        public static string EncodingNotSupported = "The required encoding (UTF-16LE) is not supported on this system";
        #endregion
    }
    #endregion

    /// <summary>
    /// Provides a stream wrapping an AESCrypt file for either encryption or decryption.
    /// The file format declare support for 2^64 bytes encrypted data, but .Net has trouble
    /// with files more than 2^63 bytes long, so this module 'only' supports 2^63 bytes
    /// (long vs ulong).
    /// </summary>
    public class SharpAESCrypt : Stream
    {
        #region Shared constant values
        /// <summary>
        /// The header in an AESCrypt file
        /// </summary>
        private readonly byte[] MAGIC_HEADER = Encoding.UTF8.GetBytes("AES");

        /// <summary>
        /// The maximum supported file version
        /// </summary>
        public const byte MAX_FILE_VERSION = 2;

        /// <summary>
        /// The size of the block unit used by the algorithm in bytes
        /// </summary>
        private const int BLOCK_SIZE = 16;
        /// <summary>
        /// The size of the IV, in bytes, which is the same as the blocksize for AES
        /// </summary>
        private const int IV_SIZE = 16;
        /// <summary>
        /// The size of the key. For AES-256 that is 256/8 = 32
        /// </summary>
        private const int KEY_SIZE = 32;
        /// <summary>
        /// The size of the SHA-256 output, which matches the KEY_SIZE
        /// </summary>
        private const int HASH_SIZE = 32;
        #endregion

        #region Private instance variables
        /// <summary>
        /// The stream being encrypted or decrypted
        /// </summary>
        private Stream m_stream;
        /// <summary>
        /// The mode of operation
        /// </summary>
        private OperationMode m_mode;
        /// <summary>
        /// The cryptostream used to perform bulk encryption
        /// </summary>
        private CryptoStream m_crypto;
        /// <summary>
        /// The HMAC used for validating data
        /// </summary>
        private HMAC m_hmac;
        /// <summary>
        /// The length of the data modulus <see cref="BLOCK_SIZE"/>
        /// </summary>
        private int m_length;
        /// <summary>
        /// The setup helper instance
        /// </summary>
        private SetupHelper m_helper;
        /// <summary>
        /// The list of extensions read from or written to the stream
        /// </summary>
        private List<KeyValuePair<string, byte[]>> m_extensions;
        /// <summary>
        /// The file format version
        /// </summary>
        private byte m_version = MAX_FILE_VERSION;
        /// <summary>
        /// True if the header is written, false otherwise. Used only for encryption.
        /// </summary>
        private bool m_hasWrittenHeader = false;
        /// <summary>
        /// True if the footer has been written, false otherwise. Used only for encryption.
        /// </summary>
        private bool m_hasFlushedFinalBlock = false;
        /// <summary>
        /// The size of the payload, including padding. Used only for decryption.
        /// </summary>
        private long m_payloadLength;
        /// <summary>
        /// The number of bytes read from the encrypted stream. Used only for decryption.
        /// </summary>
        private long m_readcount;
        /// <summary>
        /// The number of padding bytes. Used only for decryption.
        /// </summary>
        private byte m_paddingSize;
        /// <summary>
        /// True if the header HMAC has been read and verified, false otherwise. Used only for decryption.
        /// </summary>
        private bool m_hasReadFooter = false;
        #endregion

        #region Private helper functions and properties
        /// <summary>
        /// Helper property to ensure that the crypto stream is initialized before being used
        /// </summary>
        private CryptoStream Crypto
        {
            get
            {
                if (m_crypto == null)
                    WriteEncryptionHeader();
                return m_crypto;
            }
        }

        /// <summary>
        /// Helper function to read and validate the header
        /// </summary>
        private void ReadEncryptionHeader(string password)
        {
            byte[] tmp = new byte[MAGIC_HEADER.Length + 2];
            if (m_stream.Read(tmp, 0, tmp.Length) != tmp.Length)
                throw new InvalidDataException(Strings.InvalidHeaderMarker);

            for (int i = 0; i < MAGIC_HEADER.Length; i++)
                if (MAGIC_HEADER[i] != tmp[i])
                    throw new InvalidDataException(Strings.InvalidHeaderMarker);

            m_version = tmp[MAGIC_HEADER.Length];
            if (m_version > MAX_FILE_VERSION)
                throw new InvalidDataException(string.Format(Strings.UnsupportedFileVersion, m_version));

            if (m_version == 0)
            {
                m_paddingSize = tmp[MAGIC_HEADER.Length + 1];
                if (m_paddingSize >= BLOCK_SIZE)
                    throw new InvalidDataException(Strings.InvalidHeaderMarker);
            }
            else if (tmp[MAGIC_HEADER.Length + 1] != 0)
                throw new InvalidDataException(Strings.InvalidReservedFieldValue);

            //Extensions are only supported in v2+
            if (m_version >= 2)
            {
                int extensionLength = 0;
                do
                {
                    byte[] tmpLength = RepeatRead(m_stream, 2);
                    extensionLength = (((int)tmpLength[0]) << 8) | (tmpLength[1]);

                    if (extensionLength != 0)
                    {
                        byte[] data = RepeatRead(m_stream, extensionLength);
                        int separatorIndex = Array.IndexOf<byte>(data, 0);
                        if (separatorIndex < 0)
                            throw new InvalidDataException(Strings.InvalidExtensionData);

                        string key = System.Text.Encoding.UTF8.GetString(data, 0, separatorIndex);
                        byte[] value = new byte[data.Length - separatorIndex - 1];
                        Array.Copy(data, separatorIndex + 1, value, 0, value.Length);

                        m_extensions.Add(new KeyValuePair<string, byte[]>(key, value));
                    }

                } while (extensionLength > 0);
            }

            byte[] iv1 = RepeatRead(m_stream, IV_SIZE);
            m_helper = new SetupHelper(m_mode, password, iv1);

            if (m_version >= 1)
            {
                byte[] hmac1 = m_helper.DecryptAESKey2(RepeatRead(m_stream, IV_SIZE + KEY_SIZE));
                byte[] hmac2 = RepeatRead(m_stream, hmac1.Length);
                for (int i = 0; i < hmac1.Length; i++)
                    if (hmac1[i] != hmac2[i])
                        throw new CryptographicException(Strings.InvalidPassword);

                m_payloadLength = m_stream.Length - m_stream.Position - (HASH_SIZE + 1);
            }
            else
            {
                m_helper.SetBulkKeyToKey1();

                m_payloadLength = m_stream.Length - m_stream.Position - HASH_SIZE;
            }

            if (m_payloadLength % BLOCK_SIZE != 0)
                throw new CryptographicException(Strings.InvalidFileLength);
        }

        /// <summary>
        /// Writes the header to the output stream and sets up the crypto stream
        /// </summary>
        private void WriteEncryptionHeader()
        {
            m_stream.Write(MAGIC_HEADER, 0, MAGIC_HEADER.Length);
            m_stream.WriteByte(m_version);
            m_stream.WriteByte(0); //Reserved or length % 16
            if (m_version >= 2)
            {
                foreach (KeyValuePair<string, byte[]> ext in m_extensions)
                    WriteExtension(ext.Key, ext.Value);
                m_stream.Write(new byte[] { 0, 0 }, 0, 2); //No more extensions
            }

            m_stream.Write(m_helper.IV1, 0, m_helper.IV1.Length);

            if (m_version == 0)
                m_helper.SetBulkKeyToKey1();
            else
            {
                //Generate and encrypt bulk key and its HMAC
                byte[] tmpKey = m_helper.EncryptAESKey2();
                m_stream.Write(tmpKey, 0, tmpKey.Length);
                tmpKey = m_helper.CalculateKeyHmac();
                m_stream.Write(tmpKey, 0, tmpKey.Length);
            }

            m_hmac = m_helper.GetHMAC();

            //Insert the HMAC before the stream to calculate the HMAC for the ciphertext
            m_crypto = new CryptoStream(new CryptoStream(new StreamHider(m_stream, 0), m_hmac, CryptoStreamMode.Write), m_helper.CreateCryptoStream(m_mode), CryptoStreamMode.Write);
            m_hasWrittenHeader = true;
        }

        /// <summary>
        /// Writes an extension to the output stream, see:
        /// http://www.aescrypt.com/aes_file_format.html
        /// </summary>
        /// <param name="identifier">The extension identifier</param>
        /// <param name="value">The data to set in the extension</param>
        private void WriteExtension(string identifier, byte[] value)
        {
            byte[] name = System.Text.Encoding.UTF8.GetBytes(identifier);
            if (value == null)
                value = new byte[0];

            uint size = (uint)(name.Length + 1 + value.Length);
            m_stream.WriteByte((byte)((size >> 8) & 0xff));
            m_stream.WriteByte((byte)(size & 0xff));
            m_stream.Write(name, 0, name.Length);
            m_stream.WriteByte(0);
            m_stream.Write(value, 0, value.Length);
        }

        #endregion

        #region Private utility classes and functions
        /// <summary>
        /// Internal helper class used to encapsulate the setup process
        /// </summary>
        private class SetupHelper : IDisposable
        {
            /// <summary>
            /// The MAC adress to use in case the network interface enumeration fails
            /// </summary>
            private static readonly byte[] DEFAULT_MAC = { 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef };

            /// <summary>
            /// The hashing algorithm used to digest data
            /// </summary>
            private const string HASH_ALGORITHM = "SHA-256";

            /// <summary>
            /// The algorithm used to encrypt and decrypt data
            /// </summary>
            private const string CRYPT_ALGORITHM = "Rijndael";

            /// <summary>
            /// The algorithm used to generate random data
            /// </summary>
            private const string RAND_ALGORITHM = "SHA1PRNG";

            /// <summary>
            /// The algorithm used to calculate the HMAC
            /// </summary>
            private const string HMAC_ALGORITHM = "HmacSHA256";

            /// <summary>
            /// The encoding scheme for the password.
            /// UTF-16 should mean UTF-16LE, but Mono rejects the full name.
            /// A check is made when using the encoding, that it is indeed UTF-16LE.
            /// </summary>
            private const string PASSWORD_ENCODING = "utf-16";

            /// <summary>
            /// The encryption instance
            /// </summary>
            private SymmetricAlgorithm m_crypt;
            /// <summary>
            /// The hash instance
            /// </summary>
            private HashAlgorithm m_hash;
            /// <summary>
            /// The random number generator instance
            /// </summary>
            private RandomNumberGenerator m_rand;
            /// <summary>
            /// The HMAC algorithm
            /// </summary>
            private HMAC m_hmac;

            /// <summary>
            /// The IV used to encrypt/decrypt the bulk key
            /// </summary>
            private byte[] m_iv1;
            /// <summary>
            /// The private key used to encrypt/decrypt the bulk key
            /// </summary>
            private byte[] m_aesKey1;
            /// <summary>
            /// The IV used to encrypt/decrypt bulk data
            /// </summary>
            private byte[] m_iv2;
            /// <summary>
            /// The key used to encrypt/decrypt bulk data
            /// </summary>
            private byte[] m_aesKey2;

            /// <summary>
            /// Initialize the setup
            /// </summary>
            /// <param name="mode">The mode to prepare for</param>
            /// <param name="password">The password used to encrypt or decrypt</param>
            /// <param name="iv">The IV used, set to null if encrypting</param>
            public SetupHelper(OperationMode mode, string password, byte[] iv)
            {
                m_crypt = SymmetricAlgorithm.Create(CRYPT_ALGORITHM);

                //Not sure how to insert this with the CRYPT_ALGORITHM string
                m_crypt.Padding = PaddingMode.None;
                m_crypt.Mode = CipherMode.CBC;

                m_hash = HashAlgorithm.Create(HASH_ALGORITHM);
                m_rand = RandomNumberGenerator.Create(/*RAND_ALGORITHM*/);
                m_hmac = HMAC.Create(HMAC_ALGORITHM);

                if (mode == OperationMode.Encrypt)
                {
                    m_iv1 = GenerateIv1();
                    m_aesKey1 = GenerateAESKey1(EncodePassword(password));
                    m_iv2 = GenerateIv2();
                    m_aesKey2 = GenerateAESKey2();
                }
                else
                {
                    m_iv1 = iv;
                    m_aesKey1 = GenerateAESKey1(EncodePassword(password));
                }
            }

            /// <summary>
            /// Encodes the password in UTF-16LE,
            /// used to fix missing support for the full encoding
            /// name under Mono. Verifies that the encoding is correct.
            /// </summary>
            /// <param name="password">The password to encode as a byte array</param>
            /// <returns>The password encoded as a byte array</returns>
            private byte[] EncodePassword(string password)
            {
                Encoding e = Encoding.GetEncoding(PASSWORD_ENCODING);

                byte[] preamb = e == null ? null : e.GetPreamble();
                if (preamb == null || preamb.Length != 2)
                    throw new SystemException(Strings.EncodingNotSupported);

                if (preamb[0] == 0xff && preamb[1] == 0xfe)
                    return e.GetBytes(password);
                else if (preamb[0] == 0xfe && preamb[1] == 0xff)
                {
                    //We have a Big Endian, convert to Little endian
                    byte[] tmp = e.GetBytes(password);
                    if (tmp.Length % 2 != 0)
                        throw new SystemException(Strings.EncodingNotSupported);

                    for (int i = 0; i < tmp.Length; i += 2)
                    {
                        byte x = tmp[i];
                        tmp[i] = tmp[i + 1];
                        tmp[i + 1] = x;
                    }

                    return tmp;
                }
                else
                    throw new SystemException(Strings.EncodingNotSupported);
            }

            /// <summary>
            /// Gets the IV used to encrypt the bulk data key
            /// </summary>
            public byte[] IV1
            {
                get { return m_iv1; }
            }


            /// <summary>
            /// Creates the iv used for encrypting the bulk key and IV.
            /// </summary>
            /// <returns>A random IV</returns>
            private byte[] GenerateIv1()
            {
                byte[] iv = new byte[IV_SIZE];
                long time = DateTime.Now.Ticks;
                byte[] mac = null;

                /**********************************************************************
                *                                                                     *
                *   NOTE: The time and mac are COMPONENTS in the random IV input.     *
                *         The IV does not require the time or mac to be random.       *
                *                                                                     *
                *         The mac and time are used to INCREASE the ENTROPY, and      *
                *         DECOUPLE the IV from the PRNG output, in case the PRNG      *
                *         has a defect (intentional or not)                           *
                *                                                                     *
                *         Please review the DigestRandomBytes method before           *
                *         INCORRECTLY ASSUMING that the IV is generated from          *
                *         time and mac inputs.                                        *
                *                                                                     *
                ***********************************************************************/

                try
                {
                    System.Net.NetworkInformation.NetworkInterface[] interfaces = System.Net.NetworkInformation.NetworkInterface.GetAllNetworkInterfaces();
                    for (int i = 0; i < interfaces.Length; i++)
                        if (i != System.Net.NetworkInformation.NetworkInterface.LoopbackInterfaceIndex)
                        {
                            mac = interfaces[i].GetPhysicalAddress().GetAddressBytes();
                            break;
                        }
                }
                catch
                {
                    //Not much to do, just go with default MAC
                }

                if (mac == null)
                    mac = DEFAULT_MAC;

                for (int i = 0; i < 8; i++)
                    iv[i] = (byte)((time >> (i * 8)) & 0xff);

                Array.Copy(mac, 0, iv, 8, Math.Min(mac.Length, iv.Length - 8));
                return DigestRandomBytes(iv, 256);
            }

            /// <summary>
            /// Generates a key based on the IV and the password.
            /// This key is used to encrypt the actual key and IV.
            /// </summary>
            /// <param name="password">The password supplied</param>
            /// <returns>The key generated</returns>
            private byte[] GenerateAESKey1(byte[] password)
            {
                if (!m_hash.CanReuseTransform)
                    throw new CryptographicException(Strings.UnsupportedHashAlgorithmReuse);
                if (!m_hash.CanTransformMultipleBlocks)
                    throw new CryptographicException(Strings.UnsupportedHashAlgorithmBlocks);

                if (KEY_SIZE < m_hash.HashSize / 8)
                    throw new CryptographicException(string.Format(Strings.UnsupportedHashAlgorithmBlocksize, KEY_SIZE, m_hash.HashSize / 8));

                byte[] key = new byte[KEY_SIZE];
                Array.Copy(m_iv1, key, m_iv1.Length);

                for (int i = 0; i < 8192; i++)
                {
                    m_hash.Initialize();
                    m_hash.TransformBlock(key, 0, key.Length, key, 0);
                    m_hash.TransformFinalBlock(password, 0, password.Length);
                    key = m_hash.Hash;
                }

                return key;
            }

            /// <summary>
            /// Generates a random IV for encrypting data
            /// </summary>
            /// <returns>A random IV</returns>
            private byte[] GenerateIv2()
            {
                m_crypt.GenerateIV();
                return DigestRandomBytes(m_crypt.IV, 256);
            }

            /// <summary>
            /// Generates a random key for encrypting data
            /// </summary>
            /// <returns></returns>
            private byte[] GenerateAESKey2()
            {
                m_crypt.GenerateKey();
                return DigestRandomBytes(m_crypt.Key, 32);
            }

            /// <summary>
            /// Encrypts the key and IV used to encrypt data with the initial key and IV.
            /// </summary>
            /// <returns>The encrypted AES Key (including IV)</returns>
            public byte[] EncryptAESKey2()
            {
                using (MemoryStream ms = new MemoryStream())
                using (CryptoStream cs = new CryptoStream(ms, m_crypt.CreateEncryptor(m_aesKey1, m_iv1), CryptoStreamMode.Write))
                {
                    cs.Write(m_iv2, 0, m_iv2.Length);
                    cs.Write(m_aesKey2, 0, m_aesKey2.Length);
                    cs.FlushFinalBlock();

                    return ms.ToArray();
                }
            }

            /// <summary>
            /// Calculates the HMAC for the encrypted key
            /// </summary>
            /// <returns>The HMAC value</returns>
            public byte[] CalculateKeyHmac()
            {
                m_hmac.Initialize();
                m_hmac.Key = m_aesKey1;
                return m_hmac.ComputeHash(EncryptAESKey2());
            }

            /// <summary>
            /// Performs repeated hashing of the data in the byte[] combined with random data.
            /// The update is performed on the input data, which is also returned.
            /// </summary>
            /// <param name="bytes">The bytes to start the digest operation with</param>
            /// <param name="repetitions">The number of repetitions to perform</param>
            /// <returns>The digested input data, which is the same array as passed in</returns>
            private byte[] DigestRandomBytes(byte[] bytes, int repetitions)
            {
                if (bytes.Length > (m_hash.HashSize / 8))
                    throw new CryptographicException(string.Format(Strings.UnsupportedHashAlgorithmBlocksize, bytes.Length, m_hash.HashSize / 8));

                if (!m_hash.CanReuseTransform)
                    throw new CryptographicException(Strings.UnsupportedHashAlgorithmReuse);
                if (!m_hash.CanTransformMultipleBlocks)
                    throw new CryptographicException(Strings.UnsupportedHashAlgorithmBlocks);

                m_hash.Initialize();
                m_hash.TransformBlock(bytes, 0, bytes.Length, bytes, 0);
                for (int i = 0; i < repetitions; i++)
                {
                    m_rand.GetBytes(bytes);
                    m_hash.TransformBlock(bytes, 0, bytes.Length, bytes, 0);
                }

                m_hash.TransformFinalBlock(bytes, 0, 0);
                Array.Copy(m_hash.Hash, bytes, bytes.Length);
                return bytes;
            }

            /// <summary>
            /// Generates the CryptoTransform element used to encrypt/decrypt the bulk data
            /// </summary>
            /// <param name="mode">The operation mode</param>
            /// <returns>An ICryptoTransform instance</returns>
            public ICryptoTransform CreateCryptoStream(OperationMode mode)
            {
                if (mode == OperationMode.Encrypt)
                    return m_crypt.CreateEncryptor(m_aesKey2, m_iv2);
                else
                    return m_crypt.CreateDecryptor(m_aesKey2, m_iv2);
            }

            /// <summary>
            /// Creates a fresh HMAC calculation algorithm
            /// </summary>
            /// <returns>An HMAC algortihm using AES Key 2</returns>
            public HMAC GetHMAC()
            {
                HMAC h = HMAC.Create(HMAC_ALGORITHM);
                h.Key = m_aesKey2;
                return h;
            }

            /// <summary>
            /// Decrypts the bulk key and IV
            /// </summary>
            /// <param name="data">The encrypted IV followed by the key</param>
            /// <returns>The HMAC value for the key</returns>
            public byte[] DecryptAESKey2(byte[] data)
            {
                using (MemoryStream ms = new MemoryStream(data))
                using (CryptoStream cs = new CryptoStream(ms, m_crypt.CreateDecryptor(m_aesKey1, m_iv1), CryptoStreamMode.Read))
                {
                    m_iv2 = RepeatRead(cs, IV_SIZE);
                    m_aesKey2 = RepeatRead(cs, KEY_SIZE);
                }

                m_hmac.Initialize();
                m_hmac.Key = m_aesKey1;
                m_hmac.TransformFinalBlock(data, 0, data.Length);
                return m_hmac.Hash;
            }

            /// <summary>
            /// Sets iv2 and aesKey2 to iv1 and aesKey1 respectively.
            /// Used only for files with version = 0
            /// </summary>
            public void SetBulkKeyToKey1()
            {
                m_iv2 = m_iv1;
                m_aesKey2 = m_aesKey1;
            }

            #region IDisposable Members

            /// <summary>
            /// Disposes all members
            /// </summary>
            public void Dispose()
            {
                if (m_crypt != null)
                {
                    if (m_aesKey1 != null)
                        Array.Clear(m_aesKey1, 0, m_aesKey1.Length);
                    if (m_iv1 != null)
                        Array.Clear(m_iv1, 0, m_iv1.Length);
                    if (m_aesKey2 != null)
                        Array.Clear(m_aesKey2, 0, m_aesKey2.Length);
                    if (m_iv2 != null)
                        Array.Clear(m_iv2, 0, m_iv2.Length);

                    m_aesKey1 = null;
                    m_iv1 = null;
                    m_aesKey2 = null;
                    m_iv2 = null;

                    m_hash = null;
                    m_hmac = null;
                    m_rand = null;
                    m_crypt = null;
                }
            }

            #endregion
        }

        /// <summary>
        /// Internal helper class, used to hide the trailing bytes from the cryptostream
        /// </summary>
        private class StreamHider : Stream
        {
            /// <summary>
            /// The wrapped stream
            /// </summary>
            private Stream m_stream;

            /// <summary>
            /// The number of bytes to hide
            /// </summary>
            private int m_hiddenByteCount;

            /// <summary>
            /// Constructs the stream wrapper to hide the desired bytes
            /// </summary>
            /// <param name="stream">The stream to wrap</param>
            /// <param name="count">The number of bytes to hide</param>
            public StreamHider(Stream stream, int count)
            {
                m_stream = stream;
                m_hiddenByteCount = count;
            }

            #region Basic Stream implementation stuff
            public override bool CanRead { get { return m_stream.CanRead; } }
            public override bool CanSeek { get { return m_stream.CanSeek; } }
            public override bool CanWrite { get { return m_stream.CanWrite; } }
            public override void Flush() { m_stream.Flush(); }
            public override long Length { get { return m_stream.Length; } }
            public override long Seek(long offset, SeekOrigin origin) { return m_stream.Seek(offset, origin); }
            public override void SetLength(long value) { m_stream.SetLength(value); }
            public override long Position { get { return m_stream.Position; } set { m_stream.Position = value; } }
            public override void Write(byte[] buffer, int offset, int count) { m_stream.Write(buffer, offset, count); }
            #endregion

            /// <summary>
            /// The overridden read function that ensures that the caller cannot see the hidden bytes
            /// </summary>
            /// <param name="buffer">The buffer to read into</param>
            /// <param name="offset">The offset into the buffer</param>
            /// <param name="count">The number of bytes to read</param>
            /// <returns>The number of bytes read</returns>
            public override int Read(byte[] buffer, int offset, int count)
            {
                long allowedCount = Math.Max(0, Math.Min(count, m_stream.Length - (m_stream.Position + m_hiddenByteCount)));
                if (allowedCount == 0)
                    return 0;
                else
                    return m_stream.Read(buffer, offset, (int)allowedCount);
            }
        }

        /// <summary>
        /// Helper function to support reading from streams that chunck data.
        /// Will keep reading a stream until <paramref name="count"/> bytes have been read.
        /// Throws an exception if the stream is exhausted before <paramref name="count"/> bytes are read.
        /// </summary>
        /// <param name="stream">The stream to read from</param>
        /// <param name="count">The number of bytes to read</param>
        /// <returns>The data read</returns>
        internal static byte[] RepeatRead(Stream stream, int count)
        {
            byte[] tmp = new byte[count];
            while (count > 0)
            {
                int r = stream.Read(tmp, tmp.Length - count, count);
                count -= r;
                if (r == 0 && count != 0)
                    throw new InvalidDataException(Strings.UnexpectedEndOfStream);
            }

            return tmp;
        }

        #endregion

        #region Public static API

        #region Default extension control variables
        /// <summary>
        /// The name inserted as the creator software in the extensions when creating output
        /// </summary>
        public static string Extension_CreatedByIdentifier = string.Format("SharpAESCrypt v{0}", System.Reflection.Assembly.GetExecutingAssembly().GetName().Version);

        /// <summary>
        /// A value indicating if the extension data should contain the creator software
        /// </summary>
        public static bool Extension_InsertCreateByIdentifier = true;

        /// <summary>
        /// A value indicating if the extensions data should contain timestamp data
        /// </summary>
        public static bool Extension_InsertTimeStamp = false;

        /// <summary>
        /// A value indicating if the extensions data should contain an empty block as suggested by the file format
        /// </summary>
        public static bool Extension_InsertPlaceholder = true;
        #endregion

        /// <summary>
        /// The file version to use when creating a new file
        /// </summary>
        public static byte DefaultFileVersion = MAX_FILE_VERSION;

        /// <summary>
        /// Encrypts a stream using the supplied password
        /// </summary>
        /// <param name="password">The password to decrypt with</param>
        /// <param name="input">The stream with unencrypted data</param>
        /// <param name="output">The encrypted output stream</param>
        public static void Encrypt(string password, Stream input, Stream output)
        {
            int a;
            byte[] buffer = new byte[1024 * 4];
            SharpAESCrypt c = new SharpAESCrypt(password, output, OperationMode.Encrypt);
            while ((a = input.Read(buffer, 0, buffer.Length)) != 0)
                c.Write(buffer, 0, a);
            c.FlushFinalBlock();
        }

        /// <summary>
        /// Decrypts a stream using the supplied password
        /// </summary>
        /// <param name="password">The password to encrypt with</param>
        /// <param name="input">The stream with encrypted data</param>
        /// <param name="output">The unencrypted output stream</param>
        public static void Decrypt(string password, Stream input, Stream output)
        {
            int a;
            byte[] buffer = new byte[1024 * 4];
            SharpAESCrypt c = new SharpAESCrypt(password, input, OperationMode.Decrypt);
            while ((a = c.Read(buffer, 0, buffer.Length)) != 0)
                output.Write(buffer, 0, a);
        }

        /// <summary>
        /// Encrypts a file using the supplied password
        /// </summary>
        /// <param name="password">The password to encrypt with</param>
        /// <param name="inputfile">The file with unencrypted data</param>
        /// <param name="outputfile">The encrypted output file</param>
        public static void Encrypt(string password, string inputfile, string outputfile)
        {
            using (FileStream infs = File.OpenRead(inputfile))
            using (FileStream outfs = File.Create(outputfile))
                Encrypt(password, infs, outfs);
        }

        /// <summary>
        /// Decrypts a file using the supplied password
        /// </summary>
        /// <param name="password">The password to decrypt with</param>
        /// <param name="inputfile">The file with encrypted data</param>
        /// <param name="outputfile">The unencrypted output file</param>
        public static void Decrypt(string password, string inputfile, string outputfile)
        {
            using (FileStream infs = File.OpenRead(inputfile))
            using (FileStream outfs = File.Create(outputfile))
                Decrypt(password, infs, outfs);
        }
        #endregion

        #region Public instance API
        /// <summary>
        /// Constructs a new AESCrypt instance, operating on the supplied stream
        /// </summary>
        /// <param name="password">The password used for encryption or decryption</param>
        /// <param name="stream">The stream to operate on, must be writeable for encryption, and readable for decryption</param>
        /// <param name="mode">The mode of operation, either OperationMode.Encrypt or OperationMode.Decrypt</param>
        public SharpAESCrypt(string password, Stream stream, OperationMode mode)
        {
            //Basic input checks
            if (stream == null)
                throw new ArgumentNullException("stream");
            if (password == null)
                throw new ArgumentNullException("password");
            if (mode != OperationMode.Encrypt && mode != OperationMode.Decrypt)
                throw new ArgumentException(Strings.InvalidOperationMode, "mode");
            if (mode == OperationMode.Encrypt && !stream.CanWrite)
                throw new ArgumentException(Strings.StreamMustBeWriteAble, "stream");
            if (mode == OperationMode.Decrypt && !stream.CanRead)
                throw new ArgumentException(Strings.StreamMustBeReadAble, "stream");

            m_mode = mode;
            m_stream = stream;
            m_extensions = new List<KeyValuePair<string, byte[]>>();

            if (mode == OperationMode.Encrypt)
            {
                this.Version = DefaultFileVersion;

                m_helper = new SetupHelper(mode, password, null);

                //Setup default extensions
                if (Extension_InsertCreateByIdentifier)
                    m_extensions.Add(new KeyValuePair<string, byte[]>("CREATED_BY", System.Text.Encoding.UTF8.GetBytes(Extension_CreatedByIdentifier)));

                if (Extension_InsertTimeStamp)
                {
                    m_extensions.Add(new KeyValuePair<string, byte[]>("CREATED_DATE", System.Text.Encoding.UTF8.GetBytes(DateTime.UtcNow.ToString("yyyy-MM-dd"))));
                    m_extensions.Add(new KeyValuePair<string, byte[]>("CREATED_TIME", System.Text.Encoding.UTF8.GetBytes(DateTime.UtcNow.ToString("hh-mm-ss"))));
                }

                if (Extension_InsertPlaceholder)
                    m_extensions.Add(new KeyValuePair<string, byte[]>(String.Empty, new byte[127])); //Suggested extension space

                //We defer creation of the cryptostream until it is needed,
                // so the caller can change version, extensions, etc.
                // before we write the header
                m_crypto = null;
            }
            else
            {
                //Read and validate
                ReadEncryptionHeader(password);

                m_hmac = m_helper.GetHMAC();

                //Insert the HMAC before the decryption so the HMAC is calculated for the ciphertext
                m_crypto = new CryptoStream(new CryptoStream(new StreamHider(m_stream, m_version == 0 ? HASH_SIZE : (HASH_SIZE + 1)), m_hmac, CryptoStreamMode.Read), m_helper.CreateCryptoStream(m_mode), CryptoStreamMode.Read);
            }
        }

        /// <summary>
        /// Gets or sets the version number.
        /// Note that this can only be set when encrypting,
        /// and must be done before encryption has started.
        /// See <value>MAX_FILE_VERSION</value> for the maximum supported version.
        /// Note that version 0 requires a seekable stream.
        /// </summary>
        public byte Version
        {
            get { return m_version; }
            set
            {
                if (m_mode == OperationMode.Decrypt)
                    throw new InvalidOperationException(Strings.VersionReadonlyForDecryption);
                if (m_mode == OperationMode.Encrypt && m_crypto != null)
                    throw new InvalidOperationException(Strings.VersionReadonly);
                if (value > MAX_FILE_VERSION)
                    throw new ArgumentOutOfRangeException(string.Format(Strings.VersionUnsupported, MAX_FILE_VERSION));
                if (value == 0 && !m_stream.CanSeek)
                    throw new InvalidOperationException(Strings.StreamMustSupportSeeking);

                m_version = value;
            }
        }

        /// <summary>
        /// Provides access to the extensions found in the file.
        /// This collection cannot be updated when decrypting,
        /// nor after the encryption has started.
        /// </summary>
        public IList<KeyValuePair<string, byte[]>> Extensions
        {
            get
            {
                if (m_mode == OperationMode.Decrypt || (m_mode == OperationMode.Encrypt && m_crypto != null))
                    return m_extensions.AsReadOnly();
                else
                    return m_extensions;
            }
        }

        #region Basic stream implementation stuff, all mapped directly to the cryptostream
        /// <summary>
        /// Gets a value indicating whether this instance can read.
        /// </summary>
        /// <value><c>true</c> if this instance can read; otherwise, <c>false</c>.</value>
        public override bool CanRead { get { return Crypto.CanRead; } }
        /// <summary>
        /// Gets a value indicating whether this instance can seek.
        /// </summary>
        /// <value><c>true</c> if this instance can seek; otherwise, <c>false</c>.</value>
        public override bool CanSeek { get { return Crypto.CanSeek; } }
        /// <summary>
        /// Gets a value indicating whether this instance can write.
        /// </summary>
        /// <value><c>true</c> if this instance can write; otherwise, <c>false</c>.</value>
        public override bool CanWrite { get { return Crypto.CanWrite; } }
        /// <Docs>An I/O error occurs.</Docs>
        /// <summary>
        /// Flush this instance.
        /// </summary>
        public override void Flush() { Crypto.Flush(); }
        /// <summary>
        /// Gets the length.
        /// </summary>
        /// <value>The length.</value>
        public override long Length { get { return Crypto.Length; } }
        /// <summary>
        /// Gets or sets the position.
        /// </summary>
        /// <value>The position.</value>
        public override long Position
        {
            get { return Crypto.Position; }
            set { Crypto.Position = value; }
        }
        /// <Docs>The stream does not support seeking, such as if the stream is constructed from a pipe or console output.</Docs>
        /// <exception cref="T:System.IO.IOException">An I/O error has occurred.</exception>
        /// <attribution license="cc4" from="Microsoft" modified="false"></attribution>
        /// <see cref="P:System.IO.Stream.CanSeek"></see>
        /// <summary>
        /// Seek the specified offset and origin.
        /// </summary>
        /// <param name="offset">Offset.</param>
        /// <param name="origin">Origin.</param>
        public override long Seek(long offset, System.IO.SeekOrigin origin) { return Crypto.Seek(offset, origin); }
        /// <Docs>The stream does not support both writing and seeking, such as if the stream is constructed from a pipe or
        /// console output.</Docs>
        /// <exception cref="T:System.IO.IOException">An I/O error occurred.</exception>
        /// <attribution license="cc4" from="Microsoft" modified="false"></attribution>
        /// <para>A stream must support both writing and seeking for SetLength to work.</para>
        /// <see cref="P:System.IO.Stream.CanWrite"></see>
        /// <see cref="P:System.IO.Stream.CanSeek"></see>
        /// <summary>
        /// Sets the length.
        /// </summary>
        /// <param name="value">Value.</param>
        public override void SetLength(long value) { Crypto.SetLength(value); }
        #endregion

        /// <summary>
        /// Reads unencrypted data from the underlying stream
        /// </summary>
        /// <param name="buffer">The buffer to read data into</param>
        /// <param name="offset">The offset into the buffer</param>
        /// <param name="count">The number of bytes to read</param>
        /// <returns>The number of bytes read</returns>
        public override int Read(byte[] buffer, int offset, int count)
        {
            if (m_mode != OperationMode.Decrypt)
                throw new InvalidOperationException(Strings.CannotReadWhileEncrypting);

            if (m_hasReadFooter)
                return 0;

            count = Crypto.Read(buffer, offset, count);

            //TODO: If the cryptostream supporting seeking in future versions of .Net,
            // this counter system does not work
            m_readcount += count;
            m_length = (m_length + count) % BLOCK_SIZE;

            if (!m_hasReadFooter && m_readcount == m_payloadLength)
            {
                m_hasReadFooter = true;

                //Verify the data
                if (m_version >= 1)
                {
                    int l = m_stream.ReadByte();
                    if (l < 0)
                        throw new InvalidDataException(Strings.UnexpectedEndOfStream);
                    m_paddingSize = (byte)l;
                    if (m_paddingSize > BLOCK_SIZE)
                        throw new InvalidDataException(Strings.InvalidFileLength);
                }

                if (m_paddingSize > 0)
                    count -= (BLOCK_SIZE - m_paddingSize);

                if (m_length % BLOCK_SIZE != 0 || m_readcount % BLOCK_SIZE != 0)
                    throw new InvalidDataException(Strings.InvalidFileLength);

                //Required because we want to read the hash,
                // so FlushFinalBlock need to be called.
                //We cannot call FlushFinalBlock directly because it may
                // have been called by the read operation.
                //The StreamHider makes sure that the underlying stream
                // is not closed
                Crypto.Close();

                byte[] hmac1 = m_hmac.Hash;
                byte[] hmac2 = RepeatRead(m_stream, hmac1.Length);
                for (int i = 0; i < hmac1.Length; i++)
                    if (hmac1[i] != hmac2[i])
                        throw new InvalidDataException(m_version == 0 ? Strings.DataHMACMismatch_v0 : Strings.DataHMACMismatch);
            }

            return count;
        }

        /// <summary>
        /// Writes unencrypted data into an encrypted stream
        /// </summary>
        /// <param name="buffer">The data to write</param>
        /// <param name="offset">The offset into the buffer</param>
        /// <param name="count">The number of bytes to write</param>
        public override void Write(byte[] buffer, int offset, int count)
        {
            if (m_mode != OperationMode.Encrypt)
                throw new InvalidOperationException(Strings.CannotWriteWhileDecrypting);

            m_length = (m_length + count) % BLOCK_SIZE;
            Crypto.Write(buffer, offset, count);
        }

        /// <summary>
        /// Flushes any remaining data to the stream
        /// </summary>
        public void FlushFinalBlock()
        {
            if (!m_hasFlushedFinalBlock)
            {
                if (m_mode == OperationMode.Encrypt)
                {
                    if (!m_hasWrittenHeader)
                        WriteEncryptionHeader();

                    byte lastLen = (byte)(m_length %= BLOCK_SIZE);

                    //Apply PaddingMode.PKCS7 manually, the original AES crypt uses non-standard padding
                    if (lastLen != 0)
                    {
                        byte[] padding = new byte[BLOCK_SIZE - lastLen];
                        for (int i = 0; i < padding.Length; i++)
                            padding[i] = (byte)padding.Length;
                        Write(padding, 0, padding.Length);
                    }

                    //Not required without padding, but throws exception if the stream is used incorrectly
                    Crypto.FlushFinalBlock();
                    //The StreamHider makes sure the underlying stream is not closed.
                    Crypto.Close();

                    byte[] hmac = m_hmac.Hash;

                    if (m_version == 0)
                    {
                        m_stream.Write(hmac, 0, hmac.Length);
                        long pos = m_stream.Position;
                        m_stream.Seek(MAGIC_HEADER.Length + 1, SeekOrigin.Begin);
                        m_stream.WriteByte(lastLen);
                        m_stream.Seek(pos, SeekOrigin.Begin);
                        m_stream.Flush();
                    }
                    else
                    {
                        m_stream.WriteByte(lastLen);
                        m_stream.Write(hmac, 0, hmac.Length);
                        m_stream.Flush();
                    }
                }
                m_hasFlushedFinalBlock = true;
            }
        }

        /// <summary>
        /// Releases all resources used by the instance, and flushes any data currently held, into the stream
        /// </summary>
        protected override void Dispose(bool disposing)
        {
            base.Dispose(disposing);

            if (disposing)
            {
                if (m_mode == OperationMode.Encrypt && !m_hasFlushedFinalBlock)
                    FlushFinalBlock();

                if (m_crypto != null)
                    m_crypto.Dispose();
                m_crypto = null;

                if (m_stream != null)
                    m_stream.Dispose();
                m_stream = null;
                m_extensions = null;
                if (m_helper != null)
                    m_helper.Dispose();
                m_helper = null;
                m_hmac = null;
            }
        }

        #endregion



        #region Unittest code
#if DEBUG
        /// <summary>
        /// Performs a unittest to ensure that the program performs as expected
        /// </summary>
        private static void Unittest()
        {
            const int MIN_SIZE = 1024 * 5;
            const int MAX_SIZE = 1024 * 1024 * 100; //100mb
            const int REPETIONS = 1000;

            bool allpass = true;

            Random rnd = new Random();
            Console.WriteLine("Running unittest");

            //Test each supported version
            for (byte v = 0; v <= MAX_FILE_VERSION; v++)
            {
                SharpAESCrypt.DefaultFileVersion = v;

                //Test boundary 0 and around the block/keysize margins
                for (int i = 0; i < MIN_SIZE; i++)
                    using (MemoryStream ms = new MemoryStream())
                    {
                        byte[] tmp = new byte[i];
                        rnd.NextBytes(tmp);
                        ms.Write(tmp, 0, tmp.Length);
                        allpass &= Unittest(string.Format("Testing version {0} with length = {1} => ", v, ms.Length), ms);
                    }
            }

            SharpAESCrypt.DefaultFileVersion = MAX_FILE_VERSION;
            Console.WriteLine(string.Format("Initial tests complete, running bulk tests with v{0}", SharpAESCrypt.DefaultFileVersion));

            for (int i = 0; i < REPETIONS; i++)
            {
                using (MemoryStream ms = new MemoryStream())
                {
                    byte[] tmp = new byte[rnd.Next(MIN_SIZE, MAX_SIZE)];
                    rnd.NextBytes(tmp);
                    ms.Write(tmp, 0, tmp.Length);
                    allpass |= Unittest(string.Format("Testing bulk {0} of {1} with length = {2} => ", i, REPETIONS, ms.Length), ms);
                }
            }

            if (allpass)
            {
                Console.WriteLine();
                Console.WriteLine();
                Console.WriteLine("**** All unittests passed ****");
                Console.WriteLine();
            }
        }

        /// <summary>
        /// Helper function to
        /// </summary>
        /// <param name="message">A message printed to the console</param>
        /// <param name="input">The stream to test with</param>
        private static bool Unittest(string message, MemoryStream input)
        {
            Console.Write(message);

            const string PASSWORD_CHARS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!\"#Â¤%&/()=?`*'^Â¨-_.:,;<>|";
            const int MIN_LEN = 1;
            const int MAX_LEN = 25;

            try
            {
                Random rnd = new Random();
                char[] pwdchars = new char[rnd.Next(MIN_LEN, MAX_LEN)];
                for (int i = 0; i < pwdchars.Length; i++)
                    pwdchars[i] = PASSWORD_CHARS[rnd.Next(0, PASSWORD_CHARS.Length)];

                input.Position = 0;

                using (MemoryStream enc = new MemoryStream())
                using (MemoryStream dec = new MemoryStream())
                {
                    Encrypt(new string(pwdchars), input, enc);
                    enc.Position = 0;
                    Decrypt(new string(pwdchars), enc, dec);

                    dec.Position = 0;
                    input.Position = 0;

                    if (dec.Length != input.Length)
                        throw new Exception(string.Format("Length differ {0} vs {1}", dec.Length, input.Length));

                    for (int i = 0; i < dec.Length; i++)
                        if (dec.ReadByte() != input.ReadByte())
                            throw new Exception(string.Format("Streams differ at byte {0}", i));
                }
            }
            catch (Exception ex)
            {
                Console.WriteLine("FAILED: " + ex.Message);
                return false;
            }

            Console.WriteLine("OK!");
            return true;
        }
#endif
        #endregion
    }
}
